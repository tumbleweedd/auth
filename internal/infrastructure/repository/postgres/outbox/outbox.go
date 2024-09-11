package outbox

import (
	"context"
	"fmt"
	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/jmoiron/sqlx"
	outboxEntity "github.com/tumbleweedd/svc/auth_service/internal/domain/entity/outbox"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
)

type OutboxRepository struct {
	store *sqlx.DB
	log   logger.Logger

	txGetter *trmsqlx.CtxGetter
}

func NewOutboxRepository(
	store *sqlx.DB,
	log logger.Logger,
	txGetter *trmsqlx.CtxGetter,
) *OutboxRepository {
	return &OutboxRepository{
		store:    store,
		log:      log,
		txGetter: txGetter,
	}
}

func (r *OutboxRepository) Create(ctx context.Context, outbox *outboxEntity.Event) error {
	const op = "OutboxRepository.Create"

	_, err := r.txGetter.DefaultTrOrDB(ctx, r.store).ExecContext(
		ctx,
		queryCreateOutbox,
		outbox.ID,
		outbox.Type,
		outbox.Payload,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OutboxRepository) Delete(ctx context.Context, userUUIDs []string) error {
	const op = "OutboxRepository.Delete"

	_, err := r.txGetter.DefaultTrOrDB(ctx, r.store).ExecContext(
		ctx,
		queryDeleteOutbox,
		userUUIDs,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *OutboxRepository) GetEvents(ctx context.Context) ([]*outboxEntity.Event, error) {
	const op = "OutboxRepository.GetEvents"

	events := make([]*outboxEntity.Event, 0)

	err := r.store.SelectContext(ctx, &events, queryGetOutbox)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}
