package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/pprof"
)

type handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)

	GetAllOrders(c *gin.Context)
	SetOrder(c *gin.Context)

	GetBalance(c *gin.Context)
	WithdrawBalance(c *gin.Context)

	Withdrawals(c *gin.Context)
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
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.CheckMethod())

	// api
	user := router.Group("/api/user")
	{
		// routes with auth
		mainPath := user.Group("/", middleware.CheckContentTypeText(), middleware.CheckAuth(), middleware.AddUserInfoToCtx())
		mainPath.GET("orders", handler.GetAllOrders)
		mainPath.POST("orders", handler.SetOrder)
		mainPath.GET("withdrawals", handler.Withdrawals)

		balancePath := user.Group("/balance", middleware.CheckAuth(), middleware.AddUserInfoToCtx())
		balancePath.GET("/", handler.GetBalance)
		balancePath.POST("withdraw", handler.WithdrawBalance)

		// routes without auth
		nonAuth := user.Group("/", middleware.CheckContentTypeJSON())
		nonAuth.POST("register", handler.Register)
		nonAuth.POST("login", handler.Login)
	}

	// add pprof
	pprof.Register(router)

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
