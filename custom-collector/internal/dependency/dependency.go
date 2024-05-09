package dependency

import (
	"github.com/genvmoroz/custom-collector/internal/core/autocleanup"
	"github.com/genvmoroz/custom-collector/internal/http"
	"github.com/genvmoroz/custom-collector/internal/repository/stats"
	"github.com/genvmoroz/custom-collector/internal/repository/timegen"
	"github.com/samber/do"
)

type Dependency struct {
	autoCleanupTask *autocleanup.Task
	httpServer      *http.Server
}

func Build() Dependency {
	injector := do.DefaultInjector

	do.ProvideValue(injector, timegen.NewTimeGenerator())
	do.ProvideValue(injector, stats.NewRepo())

	do.Provide(injector, NewConfig)
	do.Provide(injector, NewLogger)
	do.Provide(injector, NewMemStore)
	do.Provide(injector, NewAutoCleanup)
	do.Provide(injector, NewService)
	do.Provide(injector, NewHTTPServer)

	return Dependency{
		autoCleanupTask: do.MustInvoke[*autocleanup.Task](injector),
		httpServer:      do.MustInvoke[*http.Server](injector),
	}
}

func (d *Dependency) AutoCleanupTask() *autocleanup.Task {
	return d.autoCleanupTask
}

func (d *Dependency) HTTPServer() *http.Server {
	return d.httpServer
}
