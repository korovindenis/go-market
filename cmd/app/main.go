package main

import (
	"context"
	"log"
	"os"

	"github.com/korovindenis/go-market/internal/adapters/auth"
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
	config, err := config.New()
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init logger
	logger, err := logger.New(config)
	if err != nil {
		log.Println(err)
		os.Exit(ExitWithError)
	}

	// init storage
	storage, err := storage.New(config)
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

	// init auth methods
	auth, err := auth.New(config)
	if err != nil {
		logger.Error("init auth", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// init handlers
	handler, err := handler.New(config, usecases, auth)
	if err != nil {
		logger.Error("init handler", zap.Error(err))
		os.Exit(ExitWithError)
	}

	// cancel the context when main() is terminated
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := server.Run(ctx, config, handler); err != nil {
		logger.Error("run web server", zap.Error(err))
		os.Exit(ExitWithError)
	}

	os.Exit(ExitSucces)
}
