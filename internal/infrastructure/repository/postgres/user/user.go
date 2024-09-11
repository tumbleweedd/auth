package user

import (
	"context"
	"fmt"
	userEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/user"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db       *sqlx.DB
	log      logger.Logger
	txGetter *trmsqlx.CtxGetter
}

func NewUserRepository(
	db *sqlx.DB,
	log logger.Logger,
	txGetter *trmsqlx.CtxGetter,
) *Repository {
	return &Repository{
		db:       db,
		log:      log,
		txGetter: txGetter,
	}
}

func (ur *Repository) Create(ctx context.Context, user *userEntity.User) (err error) {
	const op = "user.repository.Create"

	_, err = ur.txGetter.DefaultTrOrDB(ctx, ur.db).ExecContext(
		ctx,
		queryCreateUser,
		user.UUID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.Activated,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (ur *Repository) List(ctx context.Context) (usersDomain []*userEntity.User, err error) {
	const op = "user.repository.List"

	var users Users

	err = ur.txGetter.DefaultTrOrDB(ctx, ur.db).SelectContext(
		ctx,
		&users,
		queryListUsers,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users.ToDomain(), nil
}

func (ur *Repository) Get(ctx context.Context, uuid string) (*userEntity.User, error) {
	const op = "user.repository.Get"

	var user User

	err := ur.txGetter.DefaultTrOrDB(ctx, ur.db).GetContext(
		ctx,
		&user,
		queryGetUser,
		uuid,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user.ToDomain(), nil
}

func (ur *Repository) Delete(ctx context.Context, uuid string) error {
	const op = "user.repository.Delete"

	_, err := ur.txGetter.DefaultTrOrDB(ctx, ur.db).ExecContext(
		ctx,
		queryDeleteUser,
		uuid,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
