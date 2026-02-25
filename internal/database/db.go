package database

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is a wrapper around pgxpool.Pool for dependency injection
type Pool struct {
	*pgxpool.Pool
}

// New creates a new database pool
func New(ctx context.Context, connString string) (*Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Set connection pool settings
	config.MaxConns = 25
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Pool{pool}, nil
}

// BeginTx starts a transaction
func (p *Pool) BeginTx(ctx context.Context) (interface{}, error) {
	return p.Begin(ctx)
}

var (
	// Singleton instance
	instance *Pool
	once     sync.Once
)

// GetInstance returns the singleton database instance
func GetInstance() *Pool {
	return instance
}

// SetInstance sets the singleton database instance (for testing or initialization)
func SetInstance(p *Pool) {
	once.Do(func() {
		instance = p
	})
}
