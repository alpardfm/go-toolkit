// Package sql provides database connection management with leader/follower
// topology support, graceful shutdown, and automatic retry on connection failure.
//
// It wraps jmoiron/sqlx and supports PostgreSQL and MySQL drivers. The package
// tracks in-flight queries via the Command.Done() method, enabling graceful
// shutdown that waits for active operations to complete before closing connections.
//
// Usage:
//
//	db, err := sql.Init(cfg, logger)
//	cmd := db.Leader()
//	defer cmd.Done()
//	rows, err := cmd.Query(ctx, "get-users", "SELECT * FROM users")
package sql
