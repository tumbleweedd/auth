package service

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/tumbleweedd/svc/auth_service/internal/config"
	outboxEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/outbox"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"time"
)

type outBoxGetter interface {
	GetEvents(ctx context.Context) ([]*outboxEntity.Event, error)
}

type outBoxRemover interface {
	Delete(ctx context.Context, userUUIDs []string) error
}

type EventSenderService struct {
	outBoxRemover outBoxRemover
	outBoxGetter  outBoxGetter
	kafkaConfig   *config.KafkaConfig
	producer      sarama.SyncProducer
	log           logger.Logger
}

func New(
	outBoxRemover outBoxRemover,
	outBoxGetter outBoxGetter,
	kafkaConfig *config.KafkaConfig,
	producer sarama.SyncProducer,
	log logger.Logger,
) *EventSenderService {
	return &EventSenderService{
		outBoxRemover: outBoxRemover,
		outBoxGetter:  outBoxGetter,
		kafkaConfig:   kafkaConfig,
		producer:      producer,
		log:           log,
	}
}

func (s *EventSenderService) StartProcessingEvents(ctx context.Context, handlePeriod time.Duration) {
	const op = "service.eventSender.StartProcessingEvents"

	ticker := time.NewTicker(handlePeriod)
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.log.Info("stop processing events", fmt.Sprintf("op %s", op))

				return
			case <-ticker.C:
				if err := s.processEvents(ctx); err != nil {
					s.log.Error("error processing events", fmt.Sprintf("op %s: %s", op, err.Error()))
					continue
				}
			}
		}
	}()
}

func (s *EventSenderService) processEvents(ctx context.Context) error {
	const op = "service.eventSender.processEvents"

	events, err := s.outBoxGetter.GetEvents(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	saramaMessages := make([]*sarama.ProducerMessage, 0, len(events))
	processedMessagesIDs := make([]string, 0, len(events))

	for _, event := range events {
		saramaMessages = append(saramaMessages, &sarama.ProducerMessage{
			Topic: s.kafkaConfig.UserEventTopic,
			Value: sarama.ByteEncoder(event.Payload),
		})

		processedMessagesIDs = append(processedMessagesIDs, event.ID.String())
	}

	if err = s.producer.SendMessages(saramaMessages); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = s.outBoxRemover.Delete(ctx, processedMessagesIDs); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
