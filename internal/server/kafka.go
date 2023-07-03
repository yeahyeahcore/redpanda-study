package server

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type DepsKafka struct {
	Logger  *zap.Logger
	Brokers []string
}

type Kafka struct {
	logger *zap.Logger
	client *kadm.Client
}

func NewKafka(deps *DepsKafka) (*Kafka, error) {
	clientKGO, err := kgo.NewClient(kgo.SeedBrokers(deps.Brokers...))
	if err != nil {
		return nil, err
	}

	return &Kafka{
		client: kadm.NewClient(clientKGO),
		logger: deps.Logger,
	}, nil
}

func (receiver *Kafka) TopicExists(ctx context.Context, topic string) bool {
	topicsMetadata, err := receiver.client.ListTopics(ctx)
	if err != nil {
		receiver.logger.Error("failed to get list of topics on <TopicExists> of <AdminService>")
		return false
	}

	for _, metadata := range topicsMetadata {
		if metadata.Topic == topic {
			return true
		}
	}

	return false
}

func (receiver *Kafka) CreateTopic(ctx context.Context, topic string) error {
	responses, err := receiver.client.CreateTopics(ctx, 1, 1, nil, topic)
	if err != nil {
		receiver.logger.Error("failed to create topic on <CreateTopic> of <AdminService>")
		return err
	}

	for _, response := range responses {
		if response.Err != nil {
			receiver.logger.Error(fmt.Sprintf("unable to create topic '%s': %s", response.Topic, response.Err))
			return err
		}
	}

	return nil
}

func (receiver *Kafka) Close(_ context.Context) error {
	receiver.client.Close()

	return nil
}

func (receiver *Kafka) Initialize(ctx context.Context, topics []string) {
	for _, topic := range topics {
		if !receiver.TopicExists(ctx, topic) {
			if err := receiver.CreateTopic(ctx, topic); err != nil {
				receiver.logger.Error(fmt.Sprintf("failed to create topic <%s>", topic))
				continue
			}

			receiver.logger.Info(fmt.Sprintf("сreated topic <%s>", topic))
		}

		receiver.logger.Info(fmt.Sprintf("topic <%s> already exist", topic))
	}
}
