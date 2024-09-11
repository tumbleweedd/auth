package user

import (
	"context"
	userEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/user"
	"github.com/tumbleweedd/svc/auth_service/internal/infrastructure/cache/inmemory"
	"github.com/tumbleweedd/svc/auth_service/pkg/errorlist"
	"slices"
	"time"
)

type InMemoryUserRepository struct {
	cache *inmemory.Cache[string, *userEntity.User]
}

func NewInMemoryUserRepository(ctx context.Context, defaultTTL time.Duration) *InMemoryUserRepository {
	return &InMemoryUserRepository{
		cache: inmemory.New[string, *userEntity.User](ctx, defaultTTL),
	}
}

func (r *InMemoryUserRepository) Get(key string) (*userEntity.User, error) {
	value, ok :=  r.cache.Get(key)
	if !ok {
		return nil, errorlist.ErrCacheNotFound
	}

	return value, nil
}

func (r *InMemoryUserRepository) Add(key string, user *userEntity.User) error {
	r.cache.Add(key, user)

	return nil
}

func (r *InMemoryUserRepository) Delete(key string) error {
	r.cache.Delete(key)

	return nil
}

func (r *InMemoryUserRepository) GetByKeys(keys []string) ([]*userEntity.User, error) {
	var users []*userEntity.User

	r.cache.Range(func(k string, v *userEntity.User) bool {
		if slices.Contains(keys, k) {
			users = append(users, v)
		}

		return true
	})

	return users, nil
}
