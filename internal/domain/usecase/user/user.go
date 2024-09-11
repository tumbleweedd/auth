package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	outboxEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/outbox"
	"time"

	userEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/user"
)

type userService interface {
	Create(ctx context.Context, user *userEntity.User) (*userEntity.User, error)
	Get(ctx context.Context, uuid string) (*userEntity.User, error)
	GetByUUIDs(ctx context.Context, uuids []string) ([]*userEntity.User, error)
	Delete(ctx context.Context, uuid string) error
}

type outboxService interface {
	Create(ctx context.Context, event *outboxEntity.Event) error
}

type UseCase struct {
	userService   userService
	outboxService outboxService
}

func NewUseCase(
	userService userService,
	outboxService outboxService,
) *UseCase {
	return &UseCase{
		outboxService: outboxService,
		userService:   userService,
	}
}

func (u *UseCase) Create(ctx context.Context, user *userEntity.User) (string, error) {
	createdUser, err := u.userService.Create(ctx, user)
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(createdUser)
	if err != nil {
		return "", err
	}

	err = u.outboxService.Create(ctx, &outboxEntity.Event{
		ID:        uuid.New(),
		Type:      outboxEntity.EventTypeUserCreated,
		CreatedAt: time.Now(),
		Payload:   payload,
	})
	if err != nil {
		return "", err
	}

	return createdUser.UUID, nil
}

func (u *UseCase) Get(ctx context.Context, uuid string) (*userEntity.User, error) {
	const op = "usecase.user.Get"

	user, err := u.userService.Get(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *UseCase) Delete(ctx context.Context, uuid string) error {
	const op = "usecase.user.Delete"

	err := u.userService.Delete(ctx, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UseCase) GetByUUIDs(ctx context.Context, uuids []string) ([]*userEntity.User, error) {
	const op = "usecase.user.GetByUUIDs"

	users, err := u.userService.GetByUUIDs(ctx, uuids)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}
