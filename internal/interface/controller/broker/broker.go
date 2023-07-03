package broker

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yeahyeahcore/redpanda-study/internal/models"
	"github.com/yeahyeahcore/redpanda-study/pkg/echotools"
	"go.uber.org/zap"
)

type tariffBrokerService interface {
	Send(context.Context, *models.Tariff) error
}

type Deps struct {
	Logger              *zap.Logger
	TariffBrokerService tariffBrokerService
}

type Controller struct {
	logger              *zap.Logger
	tariffBrokerService tariffBrokerService
}

func New(deps Deps) *Controller {
	return &Controller{
		logger:              deps.Logger,
		tariffBrokerService: deps.TariffBrokerService,
	}
}

func (receiver *Controller) Send(ctx echo.Context) error {
	request, err := echotools.Bind[models.Tariff](ctx)
	if err != nil {
		receiver.logger.Error("failed to parse request body on <Send> of <BrokerController>", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := receiver.tariffBrokerService.Send(ctx.Request().Context(), request); err != nil {
		receiver.logger.Error("failed to send privileges on <Send> of <BrokerController>", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
