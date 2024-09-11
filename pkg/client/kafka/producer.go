package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
)

func NewProducer(port string, log logger.Logger) sarama.SyncProducer {
	const op = "producer.producer.NewProducer"

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{fmt.Sprintf("localhost:%s", port)}, cfg)
	if err != nil {
		log.Error(op, fmt.Sprintf("failed to start producer: %s", err.Error()))
	}

	return producer
}
