package admin

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

type Deps struct {
	Logger  *zap.Logger
	Brokers []string
}

type Kafka struct {
	logger *zap.Logger
	client *kadm.Client
}

func New(deps *Deps) (*Kafka, error) {
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

func (receiver *Kafka) CreateTopic(ctx context.Context, topic string) bool {
	responses, err := receiver.client.CreateTopics(ctx, 1, 1, nil, topic)
	if err != nil {
		receiver.logger.Error("failed to create topic on <CreateTopic> of <AdminService>")
		return false
	}

	for _, response := range responses {
		if response.Err != nil {
			receiver.logger.Error(fmt.Sprintf("Unable to create topic '%s': %s", response.Topic, response.Err))
			return false
		}

		receiver.logger.Info(fmt.Sprintf("Created topic '%s'\n", response.Topic))
	}

	return true
}

func (receiver *Kafka) Close() {
	receiver.client.Close()
}
