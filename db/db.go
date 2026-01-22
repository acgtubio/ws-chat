package db

import (
	"context"
	"fmt"

	"github.com/acgtubio/ws-chat/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewPostgresPool(ctx context.Context, config config.Config, logger *zap.SugaredLogger) (*pgxpool.Pool, error) {
	// Connect to admin postgres.
	adminConnString := fmt.Sprintf(
		"postgres://%s:%s@%s",
		config.Postgres.PostgresUsername,
		config.Postgres.PostgresPassword,
		config.Postgres.PostgresHost,
	)
	adminConn, err := pgx.Connect(ctx, adminConnString)
	if err != nil {
		logger.Errorw("Error checking initial database.",
			"error", err,
		)
		return nil, err
	}

	var exists bool
	err = adminConn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", config.Postgres.PostgresDBName).Scan(&exists)
	if err != nil {
		logger.Errorw("Error checking initial database.",
			"error", err,
		)
		return nil, err
	}

	defer adminConn.Close(ctx)

	// Creating new database.
	if !exists {
		logger.Infow(
			"Database does not exist. Creating new database",
			"dbName", config.Postgres.PostgresDBName,
		)
		createDefaultDatabase(ctx, adminConn, config.Postgres.PostgresDBName, logger)
	}

	// Create Pool
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		config.Postgres.PostgresUsername,
		config.Postgres.PostgresPassword,
		config.Postgres.PostgresHost,
		config.Postgres.PostgresDBName,
	)
	dbpool, err := pgxpool.New(ctx, connString)

	if err != nil {
		logger.Errorw("Error creating postgres pool.",
			"error", err,
		)
		return nil, err
	}

	createDefaultChatTable := `
	CREATE TABLE IF NOT EXISTS chat_history (
		id UUID PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		message STRING,
		create_date TIMESTAMPTZ DEFAULT NOW()
	)
	`

	_, err = dbpool.Query(ctx, createDefaultChatTable)
	if err != nil {
		logger.Errorw("Error creating default table.",
			"error", err,
		)
		return nil, err
	}

	return dbpool, nil
}

func createDefaultDatabase(ctx context.Context, conn *pgx.Conn, dbName string, logger *zap.SugaredLogger) error {
	_, err := conn.Query(ctx, fmt.Sprintf("CREATE DATABASE %s;", dbName))

	if err != nil {
		logger.Errorw(
			"Error creating database.",
			"error", err,
		)
	}

	return err
}
