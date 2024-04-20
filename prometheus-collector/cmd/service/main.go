package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/genvmoroz/win-stats-prometheus-collector/internal/dependency"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := start(context.Background()); err != nil {
		log.Fatalf("service crashed: %v", err)
	}
}

func start(ctx context.Context) (err error) {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	defer func() {
		if r := recover(); r != nil {
			log.Println("service panicked during bootstrapping:", r)

			err = fmt.Errorf("service panicked during bootstrapping: %v \nstack: %s", r, debug.Stack())
			return
		}
	}()

	deps := dependency.MustBuild(ctx)

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return deps.GetInfraServer().Run(ctx)
	})
	group.Go(func() error {
		return deps.GetCoreService().Collect(ctx)
	})

	return group.Wait()
}
