package logger

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type config interface {
	GetLogsLevel() string
}

type logger struct {
	*zap.Logger
}

func New(config config) (*logger, error) {
	var loggerZap *zap.Logger
	var once sync.Once

	// for singletone
	once.Do(func() {
		lvl, err := zap.ParseAtomicLevel(config.GetLogsLevel())
		if err != nil {
			panic(err)
		}
		config := zap.NewProductionConfig()
		config.Level = lvl

		zl, err := config.Build()
		if err != nil {
			panic(err)
		}
		defer zl.Sync()

		loggerZap = zl
	})

	return &logger{
		loggerZap,
	}, nil
}

func (l *logger) RequestLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// Processing request
		ctx.Next()

		endTime := time.Now()

		l.Logger.With(
			zap.Any("HTTP REQUEST", struct {
				METHOD  string
				URI     string
				STATUS  int
				LATENCY time.Duration
			}{
				ctx.Request.Method,
				ctx.Request.RequestURI,
				ctx.Writer.Status(),
				endTime.Sub(startTime),
			}),
		).Info("Request Logging")

		ctx.Next()
	}
}
