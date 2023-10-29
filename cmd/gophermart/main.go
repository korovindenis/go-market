package main

import (
	"context"
	"fmt"
	"log"

	"github.com/korovindenis/go-market/internal/adapters/accrual"
	"github.com/korovindenis/go-market/internal/adapters/auth"
	"github.com/korovindenis/go-market/internal/adapters/config"
	"github.com/korovindenis/go-market/internal/adapters/logger"
	storage "github.com/korovindenis/go-market/internal/adapters/storage/postgresql"
	"github.com/korovindenis/go-market/internal/domain/usecases"
	"github.com/korovindenis/go-market/internal/port/http/handler"
	"github.com/korovindenis/go-market/internal/port/http/middleware"
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
		log.Fatal(err)
	}
	fmt.Printf("%+v", config)

	// init logger
	logger, err := logger.New(config)
	if err != nil {
		log.Fatal(err)
	}

	// init storage
	storage, err := storage.New(config)
	if err != nil {
		logger.Fatal("init storage", zap.Error(err))
	}

	// init usecases
	usecases, err := usecases.New(config, storage)
	if err != nil {
		logger.Fatal("init usecases", zap.Error(err))
	}

	// init auth methods
	auth, err := auth.New(config)
	if err != nil {
		logger.Fatal("init auth", zap.Error(err))
	}

	// init handlers
	handler, err := handler.New(config, usecases, auth)
	if err != nil {
		logger.Fatal("init handler", zap.Error(err))
	}

	// init middleware
	middleware, err := middleware.New(config, auth)
	if err != nil {
		logger.Fatal("init middleware", zap.Error(err))
	}

	// cancel the context when main() is terminated
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// accrual
	accrual, err := accrual.New(config, storage)
	if err != nil {
		logger.Fatal("init accrual", zap.Error(err))
	}
	go accrual.Run(ctx)

	if err := server.Run(ctx, config, handler, middleware); err != nil {
		logger.Fatal("run web server", zap.Error(err))
	}
}
