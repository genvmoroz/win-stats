package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/genvmoroz/custom-collector/internal/dependency"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := start(context.Background()); err != nil {
		log.Fatal(err)
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

	deps := dependency.Build()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		deps.AutoCleanupTask().Start(ctx)
		return nil
	})
	group.Go(func() error {
		return deps.HTTPServer().Run(ctx)
	})

	return group.Wait()
}
