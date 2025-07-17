package db

import (
	"context"

	"github.com/pinjung/log-analyzer/internal/domain"
	"github.com/uptrace/bun"
)

type PostgresLogRepository struct {
	DB *bun.DB
}

func NewPostgresLogRepository(db *bun.DB) *PostgresLogRepository {
	return &PostgresLogRepository{DB: db}
}

func (r *PostgresLogRepository) SaveLog(l *domain.Log) error {
	ctx := context.Background()
	_, err := r.DB.NewInsert().Model(l).Exec(ctx)
	return err
}
