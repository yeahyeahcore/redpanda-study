package tariff

import (
	"context"
	"fmt"
	"time"

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

	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			receiver.logger.Info("read tariff messages closed")
			return
		case <-timer.C:
			privileges := receiver.tariffBrokerService.Read(ctx)

			for _, privilege := range privileges {
				fmt.Println(privilege)
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
