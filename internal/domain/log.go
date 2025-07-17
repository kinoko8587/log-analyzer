package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Log struct {
	bun.BaseModel `bun:"table:log_raw"`

	ID        int       `bun:"id,pk,autoincrement" json:"id"`
	UserID    string    `bun:"user_id,notnull" json:"user_id"`
	Timestamp time.Time `bun:"timestamp,notnull" json:"timestamp"`
	Level     string    `bun:"level,notnull" json:"level"`
	Message   string    `bun:"message" json:"message"`
}

type LogRepository interface {
	SaveLog(*Log) error
}
