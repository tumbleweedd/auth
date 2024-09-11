package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/google/uuid"
	userEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/user"
	"github.com/tumbleweedd/svc/auth_service/pkg/auth/crypt"
	"github.com/tumbleweedd/svc/auth_service/pkg/errorlist"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
)

type userInfoRepository interface {
	Create(ctx context.Context, user *userEntity.User) error
	List(ctx context.Context) ([]*userEntity.User, error)
	Delete(ctx context.Context, uuid string) error
	Get(ctx context.Context, uuid string) (*userEntity.User, error)
}

type userSecretRepository interface {
	Create(ctx context.Context, userUUID, passwordHash, passwordSalt string) error
	Delete(ctx context.Context, uuid string) error
}

type userHotStorageRepository interface {
	Add(uuid string, user *userEntity.User) error
	Get(uuid string) (*userEntity.User, error)
	Delete(uuid string) error
	GetByKeys(keys []string) ([]*userEntity.User, error)
}

type Service struct {
	userInfoRepository       userInfoRepository
	userSecretRepository     userSecretRepository
	userHotStorageRepository userHotStorageRepository
	txManager                *manager.Manager
	log                      logger.Logger
}

func NewService(
	userInfoRepository userInfoRepository,
	userSecretRepository userSecretRepository,
	userHotStorageRepository userHotStorageRepository,
	txManager *manager.Manager,
	log logger.Logger,
) *Service {
	return &Service{
		userInfoRepository:       userInfoRepository,
		userSecretRepository:     userSecretRepository,
		userHotStorageRepository: userHotStorageRepository,
		txManager:                txManager,
		log:                      log,
	}
}

func (s *Service) Create(ctx context.Context, user *userEntity.User) (*userEntity.User, error) {
	const op = "user.service.Create"

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		userUUID := uuid.New()

		user.UUID = userUUID.String()

		if err := s.userInfoRepository.Create(ctx, user); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		hash, salt, err := crypt.GetStrHashAndSalt(user.Password)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err = s.userSecretRepository.Create(ctx, userUUID.String(), string(hash), string(salt)); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err = s.userHotStorageRepository.Add(userUUID.String(), user); err != nil {
			s.log.Warn("%s: %w", op, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Get(ctx context.Context, uuid string) (*userEntity.User, error) {
	const op = "user.service.Get"

	user, err := s.userHotStorageRepository.Get(uuid)
	if err == nil {
		return user, nil
	} else if errors.Is(err, errorlist.ErrCacheNotFound) {
		s.log.Warn("%s: %w", op, err)
	}

	user, err = s.userInfoRepository.Get(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.userHotStorageRepository.Add(uuid, user); err != nil {
		s.log.Warn("%s: %w", op, err)
	}

	return user, nil
}

func (s *Service) Delete(ctx context.Context, uuid string) error {
	const op = "user.service.Delete"

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		err := s.userInfoRepository.Delete(ctx, uuid)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if err = s.userHotStorageRepository.Delete(uuid); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetByUUIDs(ctx context.Context, uuids []string) ([]*userEntity.User, error) {
	const op = "user.service.GetByUUIDs"

	users, err := s.userHotStorageRepository.GetByKeys(uuids)
	if err == nil {
		return users, nil
	} else if errors.Is(err, errorlist.ErrCacheNotFound) {
		s.log.Warn("%s: %w", op, errorlist.ErrCacheNotFound)
	}

	users, err = s.userInfoRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, user := range users {
		if err = s.userHotStorageRepository.Add(user.UUID, user); err != nil {
			s.log.Warn("%s: %w", op, err)
		}
	}

	return users, nil
}
