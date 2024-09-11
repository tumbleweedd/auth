package userSercret

import (
	"context"
	"fmt"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"time"
)

type Repository struct {
	db       *sqlx.DB
	log      logger.Logger
	txGetter *trmsqlx.CtxGetter
}

func NewUserSecretRepository(
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

func (r *Repository) Create(ctx context.Context, uuid, hash, salt string) error {
	const op = "Repository.Create"

	_, err := r.txGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		queryCreateUserSecret,
		uuid,
		hash,
		salt,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, userUUID string) error {
	const op = "Repository.Delete"

	_, err := r.txGetter.DefaultTrOrDB(ctx, r.db).ExecContext(
		ctx,
		queryDeleteUserSecret,
		userUUID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
