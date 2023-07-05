package initialize

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/yeahyeahcore/redpanda-study/internal/interface/broker/tariff"
	"go.uber.org/zap"
)

type BrokersDeps struct {
	Logger       *zap.Logger
	TariffClient *kgo.Client
}

type Brokers struct {
	TariffConsumer *tariff.Consumer
	TariffProducer *tariff.Producer
}

func NewBrokers(deps BrokersDeps) *Brokers {
	return &Brokers{
		TariffConsumer: tariff.NewConsumer(tariff.ConsumerDeps{
			Logger: deps.Logger,
			Client: deps.TariffClient,
		}),
		TariffProducer: tariff.NewProducer(tariff.ProducerDeps{
			Logger: deps.Logger,
			Client: deps.TariffClient,
		}),
	}
}
