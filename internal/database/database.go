package database

import (
	"context"
	"fmt"
	"time"

	"github.com/MaximBayurov/rate-limiter/internal/configuration"
	_ "github.com/jackc/pgx/stdlib" // justifying comment
	"github.com/jmoiron/sqlx"
)

func New(ctx context.Context, config configuration.DbConf) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", makeDsnFromConfig(config))
	if err != nil {
		return nil, fmt.Errorf("enable to connect to the database: %w", err)
	}

	go func() {
		<-ctx.Done()

		_ = db.Close()
	}()

	// Проверяем соединение
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("enable to ping database: %w", err)
	}

	// Устанавливаем настройки пула соединений
	db.SetMaxOpenConns(config.MaxOpenConn)
	db.SetMaxIdleConns(config.MaxIdleConn)
	db.SetConnMaxLifetime(config.MaxLifetimeConn * time.Minute) //nolint:durationcheck

	return db, nil
}
