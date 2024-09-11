package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tumbleweedd/svc/auth_service/internal/config"
	"github.com/tumbleweedd/svc/auth_service/pkg/logger"
	"time"
)

type PgDB struct {
	db *sqlx.DB
}

func NewClient(
	ctx context.Context,
	log logger.Logger,
	cfg *config.PgConfig,
) (*PgDB, error) {
	const op = "postgres.NewClient"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pgDB := &PgDB{db: db}

	if err = pgDB.pingContext(ctx, log); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pgDB, nil
}

func (pg *PgDB) GetDB() *sqlx.DB {
	return pg.db
}

func (pg *PgDB) Close() error {
	return pg.db.Close()
}

func (pg *PgDB) pingContext(ctx context.Context, log logger.Logger) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := pg.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
