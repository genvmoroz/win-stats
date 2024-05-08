package dependency

import (
	"github.com/genvmoroz/win-stats-service/internal/config"
	"github.com/genvmoroz/win-stats-service/internal/core"
	"github.com/genvmoroz/win-stats-service/internal/http"
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
