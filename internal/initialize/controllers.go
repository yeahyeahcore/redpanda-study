package initialize

import (
	"github.com/yeahyeahcore/redpanda-study/internal/interface/controller/broker"
	"go.uber.org/zap"
)

type ControllersDeps struct {
	Logger   *zap.Logger
	Services Services
}

type Controllers struct {
	Broker *broker.Controller
}

func NewControllers(deps ControllersDeps) *Controllers {
	return &Controllers{
		Broker: broker.New(broker.Deps{
			Logger:              deps.Logger,
			TariffBrokerService: deps.Services.TariffBroker,
		}),
	}
}
