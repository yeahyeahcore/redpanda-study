package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Deps struct {
	Client *kadm.Client
}

type Kafka struct {
	client *kadm.Client
}

func New(deps Deps) *Kafka {
	return &Kafka{client: deps.Client}
}

func (receiver *Kafka) TopicExists(ctx context.Context, topic string) (bool, error) {
	topicsMetadata, err := receiver.client.ListTopics(ctx)
	if err != nil {
		return false, err
	}

	for _, metadata := range topicsMetadata {
		if metadata.Topic == topic {
			return true, nil
		}
	}

	return false, nil
}

func (receiver *Kafka) CreateTopic(ctx context.Context, topic string) error {
	responses, err := receiver.client.CreateTopics(ctx, 1, 1, nil, topic)
	if err != nil {
		return err
	}

	for _, response := range responses {
		if response.Err != nil {
			return err
		}
	}

	return nil
}

// Initialize topics in brokers if not exist. Client connection closes self.
func Initialize(ctx context.Context, brokers, topics []string) error {
	lastInitializeErr := error(nil)

	kafkaAdminClient, err := kgo.NewClient(kgo.SeedBrokers(brokers...))
	if err != nil {
		return err
	}

	defer kafkaAdminClient.Close()

	kafka := New(Deps{Client: kadm.NewClient(kafkaAdminClient)})

	for _, topic := range topics {
		isExist, err := kafka.TopicExists(ctx, topic)
		if err != nil {
			lastInitializeErr = err
			continue
		}

		if isExist {
			continue
		}

		if err := kafka.CreateTopic(ctx, topic); err != nil {
			lastInitializeErr = fmt.Errorf("failed to create topic <%s> on <Initialize> of <BrokerService>", topic)
		}
	}

	return lastInitializeErr
}
