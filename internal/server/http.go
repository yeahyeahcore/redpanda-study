package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yeahyeahcore/redpanda-study/internal/config"
	"github.com/yeahyeahcore/redpanda-study/internal/initialize"
	"go.uber.org/zap"
)

type DepsHTTP struct {
	Logger *zap.Logger
}

type HTTP struct {
	logger *zap.Logger
	server *http.Server
	echo   *echo.Echo
}

func NewHTTP(deps DepsHTTP) *HTTP {
	echo := echo.New()

	echo.Use(middleware.Recover())

	return &HTTP{
		echo:   echo,
		logger: deps.Logger,
		server: &http.Server{
			Handler:        echo,
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    20 * time.Second,
			WriteTimeout:   20 * time.Second,
			IdleTimeout:    20 * time.Second,
		},
	}
}

func (receiver *HTTP) Listen(address string) error {
	receiver.server.Addr = address

	return receiver.server.ListenAndServe()
}

func (receiver *HTTP) Run(config *config.HTTP) {
	connectionString := fmt.Sprintf("%s:%s", config.Host, config.Port)
	startServerMessage := fmt.Sprintf("starting http server on %s", connectionString)

	receiver.logger.Info(startServerMessage)

	if err := receiver.Listen(connectionString); err != nil && err != http.ErrServerClosed {
		receiver.logger.Error("http listen error: ", zap.Error(err))
	}
}

func (receiver *HTTP) Stop(ctx context.Context) error {
	receiver.logger.Info("shutting down server...")

	if err := receiver.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}

func (receiver *HTTP) Register(controllers *initialize.Controllers) *HTTP {
	groupAmocrm := receiver.echo.Group("/broker")

	groupAmocrm.POST("/", controllers.Broker.Send)

	return receiver
}
