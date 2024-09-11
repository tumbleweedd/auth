package service

import (
	"context"
	outboxEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/outbox"
)

type outboxCreator interface {
	Create(ctx context.Context, userInfo *outboxEntity.Event) error
}

type Service struct {
	outboxRepository outboxCreator
}

func NewService(outboxRepository outboxCreator) *Service {
	return &Service{
		outboxRepository: outboxRepository,
	}
}

func (s *Service) Create(ctx context.Context, event *outboxEntity.Event) error {
	return s.outboxRepository.Create(ctx, event)
}
