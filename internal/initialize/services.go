package initialize

import (
	"github.com/yeahyeahcore/redpanda-study/internal/service/broker/tariff"
	"go.uber.org/zap"
)

type ServicesDeps struct {
	Logger  *zap.Logger
	Brokers Brokers
}

type Services struct {
	TariffBroker *tariff.Service
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		TariffBroker: tariff.New(tariff.Deps{
			Logger:   deps.Logger,
			Consumer: deps.Brokers.TariffConsumer,
			Producer: deps.Brokers.TariffProducer,
		}),
	}
}
