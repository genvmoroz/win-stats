package mem

import (
	"fmt"
	"sync"
	"time"

	"github.com/genvmoroz/custom-collector/internal/core"
	"github.com/hashicorp/go-memdb"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type (
	// database is a wrapper around memdb.MemDB to provide a table name for each database.
	database struct {
		tableName string
		*memdb.MemDB
	}

	// Store utilizes dynamic in-memory databases to manage the values associated with the sensor.
	// Creating all necessary tables from the outset isn't feasible since
	// the exact sensors are unknown in advance.
	// Furthermore, utilizing a single table for all sensors would result in decreased performance for the repository.
	Store struct {
		dbs    map[core.SensorID]database
		mux    sync.RWMutex
		logger logrus.FieldLogger
	}
)

func NewStore(logger logrus.FieldLogger) (*Store, error) {
	if lo.IsNil(logger) {
		return nil, fmt.Errorf("logger is nil")
	}

	repo := &Store{
		dbs:    make(map[core.SensorID]database),
		mux:    sync.RWMutex{},
		logger: logger,
	}

	return repo, nil
}

func (s *Store) StoreValue(sID core.SensorID, value core.Value) error {
	db, err := s.getOrCreateDB(sID)
	if err != nil {
		return err
	}

	if err = s.insert(db, fromCore(value)); err != nil {
		return err
	}

	s.logger.Debugf("[memstore] [sensor:%s] stored value: %+v\n", sID, value)

	return nil
}

// GetValuesForRange returns all values from the database that are within the specified time range.
// The range is inclusive, i.e. the records with the exact time will be included in the result.
func (s *Store) GetValuesForRange(sID core.SensorID, from, to time.Time) ([]core.Value, error) {
	if from.After(to) {
		return nil, fmt.Errorf("<from> time is after the <to> time")
	}

	db, ok := s.getDB(sID)
	if !ok {
		return nil, nil // no error, just no records
	}

	values, err := s.lowerBoundForRange(db, from, to)
	if err != nil {
		return nil, err
	}

	s.logger.Debugf(
		"[memstore] [sensor:%s] retrieved %d records for the range %s - %s\n",
		sID, len(values), from.Format(time.RFC3339), to.Format(time.RFC3339),
	)

	return values, nil
}

// DeleteOlderValues removes all values from the database that are older than the specified time
func (s *Store) DeleteOlderValues(t time.Time) error {
	s.mux.RLock()
	defer s.mux.RUnlock()

	for key, db := range s.dbs {
		n, err := s.deleteOlderValues(db, t)
		if err != nil {
			return fmt.Errorf("delete older values for %+v: %w", key, err)
		}

		s.logger.Debugf(
			"[memstore] deleted %d records older than %s for %+v\n",
			n, t.Format(time.RFC3339), key,
		)
	}

	return nil
}

func (s *Store) Close() error {
	return s.closeAll()
}

func (s *Store) closeAll() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	for key, db := range s.dbs {
		if err := s.closeDB(db); err != nil {
			return fmt.Errorf("close db for %+v: %w", key, err)
		}

		s.logger.Debugf("[memstore] closed db for %+v\n", key)
	}

	clear(s.dbs)

	s.logger.Debug("[memstore] repo closed")

	return nil
}

func (s *Store) closeDB(db database) error {
	tx := db.Txn(true)
	defer tx.Abort()

	_, err := tx.DeleteAll(db.tableName, "id")
	if err != nil {
		return fmt.Errorf("delete all records from table [%s]: %w", db.tableName, err)
	}

	tx.Commit()

	return nil
}

func (s *Store) insert(db database, dto value) error {
	txn := db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert(db.tableName, dto); err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	txn.Commit()

	return nil
}

func (s *Store) lowerBoundForRange(db database, from, to time.Time) ([]core.Value, error) {
	txn := db.Txn(false)
	defer txn.Abort()

	iter, err := txn.LowerBound(db.tableName, "id", from.UnixMilli())
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	var values []core.Value
	for {
		raw := iter.Next()
		if raw == nil || raw.(value).Timestamp > to.UnixMilli() {
			break
		}

		values = append(values, toCore(raw.(value)))
	}

	return values, nil
}

func (s *Store) deleteOlderValues(db database, t time.Time) (int, error) {
	txn := db.Txn(true)
	defer txn.Abort()

	iter, err := txn.ReverseLowerBound(db.tableName, "id", t.UnixMilli()-1) // -1 to exclude the specified time
	if err != nil {
		return 0, fmt.Errorf("reverse lower bound: %w", err)
	}

	n := 0
	for {
		raw := iter.Next()
		if raw == nil {
			break
		}

		if err = txn.Delete(db.tableName, raw); err != nil {
			return 0, fmt.Errorf("delete: %w", err)
		}
		n++
	}

	txn.Commit()

	return n, nil
}

func (s *Store) getOrCreateDB(key core.SensorID) (database, error) {
	db, ok := s.getDB(key)
	if !ok {
		if err := s.createDB(key); err != nil {
			return database{}, fmt.Errorf("create db for %+v: %w", key, err)
		}
		db, ok = s.getDB(key)
		if !ok {
			return database{}, fmt.Errorf("table for %+v not found", key) // should never happen
		}
	}

	return db, nil
}

func (s *Store) getDB(key core.SensorID) (database, bool) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if db, ok := s.dbs[key]; ok {
		return db, true
	}

	return database{}, false
}

func (s *Store) createDB(key core.SensorID) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, ok := s.dbs[key]; ok {
		return nil // already exists
	}

	tableName := string(key)

	memDB, err := newMemDB(tableName)
	if err != nil {
		return fmt.Errorf("new memdb: %w", err)
	}

	s.dbs[key] = database{
		tableName: tableName,
		MemDB:     memDB,
	}

	return nil
}

func newMemDB(name string) (*memdb.MemDB, error) {
	return memdb.NewMemDB(
		&memdb.DBSchema{
			Tables: map[string]*memdb.TableSchema{
				name: {
					Name: name,
					Indexes: map[string]*memdb.IndexSchema{
						"id": {
							Name:    "id",
							Unique:  true,
							Indexer: &memdb.IntFieldIndex{Field: "Timestamp"},
						},
						"value": {
							Name:    "value",
							Unique:  false,
							Indexer: &memdb.IntFieldIndex{Field: "Value"},
						},
					},
				},
			},
		},
	)
}
