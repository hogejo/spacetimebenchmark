package main

import (
	"context"
	"fmt"
	"time"

	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/schema.sql
var schemaSQL string

//go:embed sql/transfer.sql
var transferSQL string

const (
	connectionTimeout = 3 * time.Second
)

func openPool(config Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}
	poolConfig.MaxConns = int32(config.maximumConnections)
	poolConfig.MinConns = int32(config.maximumConnections)
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open db pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return pool, nil
}

func closePool(pool *pgxpool.Pool) {
	pool.Close()
}

func acquireConnections(pool *pgxpool.Pool, num int) ([]*pgxpool.Conn, error) {
	connections := make([]*pgxpool.Conn, num)
	for i := range connections {
		c, err := pool.Acquire(context.Background())
		if err != nil {
			// Release already acquired connections
			for j := 0; j < i; j++ {
				connections[j].Release()
			}
			return nil, fmt.Errorf("failed to acquire connection %d: %v", i, err)
		}
		connections[i] = c
	}
	return connections, nil
}

func releaseConnections(connections []*pgxpool.Conn) {
	for _, c := range connections {
		c.Release()
	}
}

func prepareDatabase(pool *pgxpool.Pool, config Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()
	fmt.Printf("preparing the database ...")
	for _, query := range []string{schemaSQL, transferSQL} {
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to prepare the database: %w", err)
		}
	}
	fmt.Println(" done.")
	fmt.Printf("seeding the database ...")
	if _, err := pool.Exec(ctx, "DELETE FROM balances"); err != nil {
		return fmt.Errorf("failed to clear balances: %w", err)
	}
	if _, err := pool.Exec(ctx,
		"INSERT INTO balances (id, balance) SELECT g, $2 FROM generate_series(0, $1 - 1) AS g",
		config.accounts, config.initialBalance,
	); err != nil {
		return fmt.Errorf("failed to seed balances: %w", err)
	}
	fmt.Println(" done.")
	return nil
}
