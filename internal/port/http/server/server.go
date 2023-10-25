package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

type handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)

	GetOrder(c *gin.Context)
	SetOrder(c *gin.Context)
}

type middleware interface {
	CheckMethod() gin.HandlerFunc

	CheckContentTypeJSON() gin.HandlerFunc
	CheckContentTypeText() gin.HandlerFunc

	CheckAuth() gin.HandlerFunc
	AddUserInfoToCtx() gin.HandlerFunc
}

type config interface {
	GetServerAddress() string
	GetServerMode() string
	GetServerTimeoutRead() time.Duration
	GetServerTimeoutWrite() time.Duration
	GetServerTimeoutIdle() time.Duration
	GetServerMaxHeaderBytes() int
}

func Run(ctx context.Context, config config, handler handler, middleware middleware) error {
	// init http
	gin.SetMode(config.GetServerMode())
	router := gin.Default()

	// middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CheckMethod())

	// api
	user := router.Group("/api/user")
	{
		// routes without auth
		nonAuth := user.Group("/", middleware.CheckContentTypeJSON())
		nonAuth.POST("register", handler.Register)
		nonAuth.POST("login", handler.Login)

		// routes with auth
		orders := user.Group("/", middleware.CheckContentTypeText(), middleware.CheckAuth(), middleware.AddUserInfoToCtx())
		orders.GET("orders", handler.GetOrder)
		orders.POST("orders", handler.SetOrder)
	}

	// server settings
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
