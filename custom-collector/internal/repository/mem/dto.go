package mem

import (
	"time"

	"github.com/genvmoroz/custom-collector/internal/core"
)

type value struct {
	Value     int64
	Timestamp int64
}

func toCore(dto value) core.Value {
	return core.Value{
		Value:     dto.Value,
		Timestamp: time.UnixMilli(dto.Timestamp),
	}
}

func fromCore(temp core.Value) value {
	return value{
		Value:     temp.Value,
		Timestamp: temp.Timestamp.UnixMilli(),
	}
}
