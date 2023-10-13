package main

import (
	"context"
	"log"
	"os"

	"github.com/korovindenis/go-market/internal/adapters/config"
	"github.com/korovindenis/go-market/internal/adapters/logger"
	storage "github.com/korovindenis/go-market/internal/adapters/storage/postgresql"
	"github.com/korovindenis/go-market/internal/domain/usecases"
	"github.com/korovindenis/go-market/internal/port/http/handler"
	"github.com/korovindenis/go-market/internal/port/http/server"
	"go.uber.org/zap"
)

const (
	ExitSucces = iota
	ExitWithError
)

func main() {
	// init config
	cfg, err := config.New()
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}
	// init logger
	logger, err := logger.New(cfg)
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init storage
	storage, err := storage.New(cfg)
	if err != nil {
		logger.Error("init storage", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// init usecases
	usecases, err := usecases.New(storage)
	if err != nil {
		logger.Error("init usecases", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// init handlers
	handler, err := handler.New(usecases)
	if err != nil {
		logger.Error("init handler", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// cancel the context when main() is terminated
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := server.Run(ctx, cfg, handler); err != nil {
		logger.Error("run web server", zap.Error(err))
		os.Exit(ExitWithError)
	}

	os.Exit(ExitSucces)
}
