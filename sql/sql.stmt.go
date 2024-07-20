package sql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type CommandStmt interface {
	Close() error
	Select(name string, dest interface{}, args ...interface{}) error
	Get(name string, dest interface{}, args ...interface{}) error
	QueryRow(name string, args ...interface{}) (*sqlx.Row, error)
	Query(name string, args ...interface{}) (*sqlx.Rows, error)
	Exec(name string, args ...interface{}) (sql.Result, error)
}

type commandStmt struct {
	ctx           context.Context
	name          string
	stmt          *sqlx.Stmt
	useInstrument bool
}

func initStmt(ctx context.Context, name string, stmt *sqlx.Stmt) CommandStmt {
	return &commandStmt{
		ctx:  ctx,
		name: name,
		stmt: stmt,
	}
}

func (x *commandStmt) Close() error {
	return x.stmt.Close()
}

func (x *commandStmt) Select(name string, dest interface{}, args ...interface{}) error {
	return x.stmt.SelectContext(x.ctx, dest, args...)
}

func (x *commandStmt) Get(name string, dest interface{}, args ...interface{}) error {
	return x.stmt.GetContext(x.ctx, dest, args...)
}

func (x *commandStmt) QueryRow(name string, args ...interface{}) (*sqlx.Row, error) {
	return x.stmt.QueryRowxContext(x.ctx, args...), nil
}

func (x *commandStmt) Query(name string, args ...interface{}) (*sqlx.Rows, error) {
	return x.stmt.QueryxContext(x.ctx, args...)
}

func (x *commandStmt) Exec(name string, args ...interface{}) (sql.Result, error) {
	return x.stmt.ExecContext(x.ctx, args...)
}
