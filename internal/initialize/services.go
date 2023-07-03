package initialize

import (
	"github.com/yeahyeahcore/redpanda-study/internal/service/broker/tariff"
	"go.uber.org/zap"
)

type ServicesDeps struct {
	Logger *zap.Logger
	Kafka  Kafka
}

type Services struct {
	TariffBrokerService *tariff.Service
}

func NewServices(deps ServicesDeps) *Services {
	return &Services{
		TariffBrokerService: tariff.New(tariff.Deps{
			Logger:   deps.Logger,
			Consumer: deps.Kafka.TariffConsumer,
			Producer: deps.Kafka.TariffProducer,
		}),
	}
}
