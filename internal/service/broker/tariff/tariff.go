package tariff

import (
	"context"

	"github.com/yeahyeahcore/redpanda-study/internal/models"
	"go.uber.org/zap"
)

type consumer interface {
	FetchPrivileges(context.Context) []models.Privilege
}

type producer interface {
	Send(context.Context, *models.Tariff) error
}

type Deps struct {
	Logger   *zap.Logger
	Consumer consumer
	Producer producer
}

type Service struct {
	logger   *zap.Logger
	consumer consumer
	producer producer
}

func New(deps Deps) *Service {
	return &Service{
		logger:   deps.Logger,
		consumer: deps.Consumer,
		producer: deps.Producer,
	}
}

func (receiver *Service) Send(ctx context.Context, tariff *models.Tariff) error {
	if err := receiver.producer.Send(ctx, tariff); err != nil {
		receiver.logger.Error("failed to send tariff on <Send> of <TariffBrokerService>", zap.Error(err))
		return err
	}

	return nil
}

func (receiver *Service) Read(ctx context.Context) []models.Privilege {
	return receiver.consumer.FetchPrivileges(ctx)
}
