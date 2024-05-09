package dependency

import (
	"github.com/genvmoroz/custom-collector/internal/repository/mem"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

func NewMemStore(injector *do.Injector) (*mem.Store, error) {
	logger := do.MustInvoke[logrus.FieldLogger](injector)

	return mem.NewStore(logger)
}
