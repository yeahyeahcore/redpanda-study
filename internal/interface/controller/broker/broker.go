package broker

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yeahyeahcore/redpanda-study/internal/interface/controller/broker/dto"
	"github.com/yeahyeahcore/redpanda-study/internal/models"
	"github.com/yeahyeahcore/redpanda-study/pkg/echotools"
	"go.uber.org/zap"
)

type TariffBrokerService interface {
	Send(context.Context, []models.Privilege) error
}

type Deps struct {
	Logger              *zap.Logger
	TariffBrokerService TariffBrokerService
}

type Controller struct {
	logger              *zap.Logger
	tariffBrokerService TariffBrokerService
}

func New(deps Deps) *Controller {
	return &Controller{
		logger:              deps.Logger,
		tariffBrokerService: deps.TariffBrokerService,
	}
}

func (receiver *Controller) Send(ctx echo.Context) error {
	request, err := echotools.Bind[dto.Tariff](ctx)
	if err != nil {
		receiver.logger.Error("failed to parse request body on <Send> of <BrokerController>", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := receiver.tariffBrokerService.Send(ctx.Request().Context(), request.Privileges); err != nil {
		receiver.logger.Error("failed to send privileges on <Send> of <BrokerController>", zap.Error(err))
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
