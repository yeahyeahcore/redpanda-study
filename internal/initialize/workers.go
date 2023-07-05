package initialize

import (
	"github.com/yeahyeahcore/redpanda-study/internal/worker/broker/tariff"
	"go.uber.org/zap"
)

type WorkersDeps struct {
	Logger   *zap.Logger
	Services Services
}

type Workers struct {
	TariffBroker *tariff.Worker
}

func NewWorkers(deps WorkersDeps) *Workers {
	return &Workers{
		TariffBroker: tariff.New(tariff.Deps{
			Logger:              deps.Logger,
			TariffBrokerService: deps.Services.TariffBroker,
		}),
	}
}
