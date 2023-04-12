package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/nickzhog/userapi/internal/server/config"
	"github.com/nickzhog/userapi/internal/server/repositories"
	"github.com/nickzhog/userapi/internal/server/web"
	"github.com/nickzhog/userapi/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()
	logger.Tracef("%+v", *cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		oscall := <-c
		logger.Tracef("system call:%+v", oscall)
		cancel()
	}()

	reps := repositories.GetRepositories(ctx, logger, *cfg)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		srv := web.PrepareServer(cfg.Settings.RunAddress, logger, reps)
		if err := web.Serve(ctx, logger, srv); err != nil {
			logger.Errorf("failed to serve: %s", err.Error())
		}
		wg.Done()
	}()

	wg.Wait()
}
