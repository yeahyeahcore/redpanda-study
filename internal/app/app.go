package app

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/yeahyeahcore/redpanda-study/internal/config"
	"github.com/yeahyeahcore/redpanda-study/internal/initialize"
	"github.com/yeahyeahcore/redpanda-study/internal/server"
	"github.com/yeahyeahcore/redpanda-study/internal/worker"
	"github.com/yeahyeahcore/redpanda-study/pkg/closer"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second
)

func Run(config *config.Config, logger *zap.Logger) (appErr error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			appErr = errors.New("recovered panic on application run")
			logger.Error("recovered panic on application run", zap.Any("recovered", recovered))
		}
	}()

	closer := closer.New()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	kafka, err := initialize.NewKafka(initialize.KafkaDeps{
		Logger: logger,
		Config: config.Service.Kafka,
	})
	if err != nil {
		return err
	}

	services := initialize.NewServices(initialize.ServicesDeps{
		Logger: logger,
		Kafka:  *kafka,
	})

	workers := initialize.NewWorkers(initialize.WorkersDeps{
		Logger:   logger,
		Services: *services,
	})

	controllers := initialize.NewControllers(initialize.ControllersDeps{Logger: logger})

	serverHTTP := server.NewHTTP(server.DepsHTTP{Logger: logger})

	serverKafka, kafkaErr := server.NewKafka(&server.DepsKafka{
		Logger:  logger,
		Brokers: config.Service.Kafka.Brokers,
	})
	if kafkaErr != nil {
		return kafkaErr
	}

	serverHTTP.Register(controllers)

	go serverHTTP.Run(&config.HTTP)
	go serverKafka.Initialize(ctx, []string{config.Service.Kafka.Tariff.Topic})
	go worker.Run(ctx, workers)

	closer.Add(serverHTTP.Stop)
	closer.Add(serverKafka.Close)
	closer.Add(kafka.TariffConsumer.Close)
	closer.Add(kafka.TariffProducer.Close)

	<-ctx.Done()

	logger.Info("shutting down app gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

	defer cancel()

	if err := closer.Close(shutdownCtx); err != nil {
		return fmt.Errorf("closer: %w", err)
	}

	return nil
}
