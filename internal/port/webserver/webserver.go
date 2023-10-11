package webserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	errTypeIsNotConfig = errors.New("webserver, type is no config")
)

type config interface {
	GetServerAddress() string
	GetServerMode() string
	GetServerTimeoutRead() time.Duration
	GetServerTimeoutWrite() time.Duration
	GetServerTimeoutIdle() time.Duration
	GetServerMaxHeaderBytes() int
}

func Run(ctx context.Context, cfgAny any) error {
	cfg, ok := cfgAny.(config)
	if !ok {
		return errTypeIsNotConfig
	}

	// init http
	gin.SetMode(cfg.GetServerMode())
	router := gin.Default()

	// middleware
	router.Use(gin.Recovery())

	// Define endpoint
	//router.GET("/", computerhandler.MainPageHandler)

	srv := &http.Server{
		Addr:           cfg.GetServerAddress(),
		Handler:        router,
		ReadTimeout:    cfg.GetServerTimeoutRead(),
		WriteTimeout:   cfg.GetServerTimeoutWrite(),
		IdleTimeout:    cfg.GetServerTimeoutIdle(),
		MaxHeaderBytes: cfg.GetServerMaxHeaderBytes(),
	}

	// run with graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Failed to listen and serve", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	return srv.Shutdown(ctx)
}
