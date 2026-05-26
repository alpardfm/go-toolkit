package sql

import (
	"context"
	gosql "database/sql"
	"sync"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// mockLog implements log.Interface for testing purposes.
type mockLog struct{}

func (m *mockLog) Trace(_ context.Context, _ any) {}
func (m *mockLog) Debug(_ context.Context, _ any) {}
func (m *mockLog) Info(_ context.Context, _ any)  {}
func (m *mockLog) Warn(_ context.Context, _ any)  {}
func (m *mockLog) Error(_ context.Context, _ any) {}
func (m *mockLog) Fatal(_ context.Context, _ any) {}

// mockCommand implements Command for testing purposes.
type mockCommand struct {
	closed bool
	mu     sync.Mutex
}

func (mc *mockCommand) Close() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.closed = true
	return nil
}

func (mc *mockCommand) Ping(_ context.Context) error { return nil }
func (mc *mockCommand) In(_ string, _ ...any) (string, []any, error) {
	return "", nil, nil
}
func (mc *mockCommand) Rebind(_ string) string { return "" }
func (mc *mockCommand) QueryIn(_ context.Context, _ string, _ string, _ ...any) (*sqlx.Rows, error) {
	return nil, nil
}
func (mc *mockCommand) QueryRow(_ context.Context, _ string, _ string, _ ...any) (*sqlx.Row, error) {
	return nil, nil
}
func (mc *mockCommand) Query(_ context.Context, _ string, _ string, _ ...any) (*sqlx.Rows, error) {
	return nil, nil
}
func (mc *mockCommand) NamedQuery(_ context.Context, _ string, _ string, _ any) (*sqlx.Rows, error) {
	return nil, nil
}
func (mc *mockCommand) Prepare(_ context.Context, _ string, _ string) (CommandStmt, error) {
	return nil, nil
}
func (mc *mockCommand) NamedExec(_ context.Context, _ string, _ string, _ any) (gosql.Result, error) {
	return nil, nil
}
func (mc *mockCommand) Exec(_ context.Context, _ string, _ string, _ ...any) (gosql.Result, error) {
	return nil, nil
}
func (mc *mockCommand) BeginTx(_ context.Context, _ string, _ TxOptions) (CommandTx, error) {
	return nil, nil
}
func (mc *mockCommand) Get(_ context.Context, _ string, _ string, _ any, _ ...any) error {
	return nil
}
func (mc *mockCommand) Done() {}

func newTestSQLDB() *sqlDB {
	return &sqlDB{
		endOnce:  &sync.Once{},
		leader:   &mockCommand{},
		follower: &mockCommand{},
		log:      &mockLog{},
	}
}

func TestShutdown_AllQueriesCompleteBeforeDeadline(t *testing.T) {
	db := newTestSQLDB()

	// Simulate an in-flight query that completes quickly
	db.inflight.Add(1)
	go func() {
		time.Sleep(10 * time.Millisecond)
		db.inflight.Done()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestShutdown_TimeoutWhenQueriesExceedDeadline(t *testing.T) {
	db := newTestSQLDB()

	// Simulate an in-flight query that takes too long
	db.inflight.Add(1)
	go func() {
		time.Sleep(5 * time.Second)
		db.inflight.Done()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := db.Shutdown(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "shutdown timed out")
}

func TestShutdown_DefaultTimeoutWhenNoDeadline(t *testing.T) {
	db := newTestSQLDB()

	// No in-flight queries, should complete immediately with default timeout
	err := db.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestShutdown_SetsClosingFlag(t *testing.T) {
	db := newTestSQLDB()

	err := db.Shutdown(context.Background())
	assert.NoError(t, err)

	db.shutdownMu.Lock()
	assert.True(t, db.closing)
	db.shutdownMu.Unlock()
}

func TestShutdown_ClosesConnections(t *testing.T) {
	leader := &mockCommand{}
	follower := &mockCommand{}
	db := &sqlDB{
		endOnce:  &sync.Once{},
		leader:   leader,
		follower: follower,
		log:      &mockLog{},
	}

	err := db.Shutdown(context.Background())
	assert.NoError(t, err)

	leader.mu.Lock()
	assert.True(t, leader.closed)
	leader.mu.Unlock()

	follower.mu.Lock()
	assert.True(t, follower.closed)
	follower.mu.Unlock()
}
