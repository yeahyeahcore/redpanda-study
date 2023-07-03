package initialize

import (
	"github.com/yeahyeahcore/redpanda-study/internal/config"
	"github.com/yeahyeahcore/redpanda-study/internal/service/broker/tariff"
	"go.uber.org/zap"
)

type BrokersDeps struct {
	Logger *zap.Logger
	Config config.KafkaConfiguration
}

type Brokers struct {
	TariffConsumer *tariff.Consumer
	TariffProducer *tariff.Producer
}

func NewBrokers(deps *BrokersDeps) (*Brokers, error) {
	tariffConsumer, err := tariff.NewConsumer(&tariff.ConsumerDeps{
		Logger:  deps.Logger,
		Brokers: deps.Config.Tariff.Brokers,
		Topic:   deps.Config.Tariff.Topic,
		GroupID: deps.Config.Tariff.GroupID,
	})
	if err != nil {
		return nil, err
	}

	tariffProducer, err := tariff.NewProducer(&tariff.ProducerDeps{
		Logger:  deps.Logger,
		Brokers: deps.Config.Tariff.Brokers,
		Topic:   deps.Config.Tariff.Topic,
	})
	if err != nil {
		return nil, err
	}

	return &Brokers{
		TariffConsumer: tariffConsumer,
		TariffProducer: tariffProducer,
	}, nil
}
