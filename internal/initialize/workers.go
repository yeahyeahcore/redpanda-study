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
	TariffBrokerWorker *tariff.Worker
}

func NewWorkers(deps WorkersDeps) *Workers {
	return &Workers{
		TariffBrokerWorker: tariff.New(tariff.WorkerDeps{
			Logger:              deps.Logger,
			TariffBrokerService: deps.Services.TariffBrokerService,
		}),
	}
}
