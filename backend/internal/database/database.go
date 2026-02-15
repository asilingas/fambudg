package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/asilingas/fambudg/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates a new connection pool to the database
func NewPool(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	// Build connection string
	connString := cfg.ConnectionString() + " pool_max_conns=10"

	// Parse config
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Set connection pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30
	poolConfig.HealthCheckPeriod = time.Minute

	// Create pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Database connection pool established successfully")
	return pool, nil
}
