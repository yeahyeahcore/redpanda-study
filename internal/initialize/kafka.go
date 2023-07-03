package initialize

import (
	"github.com/yeahyeahcore/redpanda-study/internal/config"
	"github.com/yeahyeahcore/redpanda-study/internal/interface/broker/tariff"
	"go.uber.org/zap"
)

type KafkaDeps struct {
	Logger *zap.Logger
	Config config.Kafka
}

type Kafka struct {
	TariffConsumer *tariff.Consumer
	TariffProducer *tariff.Producer
}

func NewKafka(deps KafkaDeps) (*Kafka, error) {
	tariffConsumer, err := tariff.NewConsumer(tariff.ConsumerDeps{
		Logger:  deps.Logger,
		Brokers: deps.Config.Brokers,
		Topic:   deps.Config.Tariff.Topic,
		GroupID: deps.Config.Tariff.GroupID,
	})
	if err != nil {
		return nil, err
	}

	tariffProducer, err := tariff.NewProducer(tariff.ProducerDeps{
		Logger:  deps.Logger,
		Brokers: deps.Config.Brokers,
		Topic:   deps.Config.Tariff.Topic,
	})
	if err != nil {
		return nil, err
	}

	return &Kafka{
		TariffConsumer: tariffConsumer,
		TariffProducer: tariffProducer,
	}, nil
}
