package dependency

import (
	"github.com/genvmoroz/custom-collector/internal/config"
	"github.com/genvmoroz/custom-collector/internal/core"
	"github.com/genvmoroz/custom-collector/internal/http"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

func NewHTTPServer(injector *do.Injector) (*http.Server, error) {
	var (
		cfg    = do.MustInvoke[config.Config](injector)
		srv    = do.MustInvoke[*core.Service](injector)
		logger = do.MustInvoke[logrus.FieldLogger](injector)
	)

	return http.NewServer(cfg.HTTPServer, srv, logger)
}
