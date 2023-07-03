package tariff

import (
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/yeahyeahcore/redpanda-study/internal/models"
	"go.uber.org/zap"
)

type tariffBrokerService interface {
	Send(context.Context, *models.Tariff) error
	Read(context.Context) []models.Privilege
}

type WorkerDeps struct {
	Logger              *zap.Logger
	TariffBrokerService tariffBrokerService
}

type Worker struct {
	logger              *zap.Logger
	tariffBrokerService tariffBrokerService
}

func New(deps WorkerDeps) *Worker {
	return &Worker{
		logger:              deps.Logger,
		tariffBrokerService: deps.TariffBrokerService,
	}
}

func (receiver *Worker) Run(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			receiver.logger.Error("accrual worker recover", zap.Any("recovered error", err))
		}
	}()

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ctx.Done():
			receiver.logger.Info("read tariff messages closed")
			return
		case <-ticker.C:
			tariffCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			privileges := receiver.tariffBrokerService.Read(tariffCtx)

			cancel()
			spew.Dump(privileges)
		}
	}
}
