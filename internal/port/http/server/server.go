package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/port/http/middleware"
)

type handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type config interface {
	GetServerAddress() string
	GetServerMode() string
	GetServerTimeoutRead() time.Duration
	GetServerTimeoutWrite() time.Duration
	GetServerTimeoutIdle() time.Duration
	GetServerMaxHeaderBytes() int
}

func Run(ctx context.Context, config config, handler handler) error {
	// init http
	gin.SetMode(config.GetServerMode())
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(middleware.CheckMethodAndContentType())
	nonAuthenticatedGroup := router.Group("/api/user")
	{
		nonAuthenticatedGroup.POST("/register/", handler.Register)
		nonAuthenticatedGroup.POST("/login/", handler.Login)
	}

	srv := &http.Server{
		Addr:           config.GetServerAddress(),
		Handler:        router,
		ReadTimeout:    config.GetServerTimeoutRead(),
		WriteTimeout:   config.GetServerTimeoutWrite(),
		IdleTimeout:    config.GetServerTimeoutIdle(),
		MaxHeaderBytes: config.GetServerMaxHeaderBytes(),
	}

	// run with graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("failed to listen and serve", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	return srv.Shutdown(ctx)
}
