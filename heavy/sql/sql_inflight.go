package sql

import (
	"context"
	"database/sql"
	"sync"

	"github.com/jmoiron/sqlx"
)

// inflightCommand wraps a Command and decrements the inflight WaitGroup
// when the caller signals completion via Done(). This enables graceful shutdown
// to wait for in-flight queries to finish.
//
// Usage: call Done() after you are finished with the Command returned by Leader() or Follower().
// If you use BeginTx, call Done() after the transaction is committed or rolled back.
type inflightCommand struct {
	cmd  Command
	wg   *sync.WaitGroup
	done bool
}

// Done signals that the caller is finished with this command.
// Must be called exactly once per Leader()/Follower() call to allow graceful shutdown.
func (ic *inflightCommand) Done() {
	if !ic.done {
		ic.done = true
		ic.wg.Done()
	}
}

func (ic *inflightCommand) Close() error {
	return ic.cmd.Close()
}

func (ic *inflightCommand) Ping(ctx context.Context) error {
	return ic.cmd.Ping(ctx)
}

func (ic *inflightCommand) In(query string, args ...any) (string, []any, error) {
	return ic.cmd.In(query, args...)
}

func (ic *inflightCommand) Rebind(query string) string {
	return ic.cmd.Rebind(query)
}

func (ic *inflightCommand) QueryIn(ctx context.Context, name string, query string, args ...any) (*sqlx.Rows, error) {
	return ic.cmd.QueryIn(ctx, name, query, args...)
}

func (ic *inflightCommand) QueryRow(ctx context.Context, name string, query string, args ...any) (*sqlx.Row, error) {
	return ic.cmd.QueryRow(ctx, name, query, args...)
}

func (ic *inflightCommand) Query(ctx context.Context, name string, query string, args ...any) (*sqlx.Rows, error) {
	return ic.cmd.Query(ctx, name, query, args...)
}

func (ic *inflightCommand) NamedQuery(ctx context.Context, name string, query string, arg any) (*sqlx.Rows, error) {
	return ic.cmd.NamedQuery(ctx, name, query, arg)
}

func (ic *inflightCommand) Prepare(ctx context.Context, name string, query string) (CommandStmt, error) {
	return ic.cmd.Prepare(ctx, name, query)
}

func (ic *inflightCommand) NamedExec(ctx context.Context, name string, query string, args any) (sql.Result, error) {
	return ic.cmd.NamedExec(ctx, name, query, args)
}

func (ic *inflightCommand) Exec(ctx context.Context, name string, query string, args ...any) (sql.Result, error) {
	return ic.cmd.Exec(ctx, name, query, args...)
}

func (ic *inflightCommand) BeginTx(ctx context.Context, name string, opts TxOptions) (CommandTx, error) {
	return ic.cmd.BeginTx(ctx, name, opts)
}

func (ic *inflightCommand) Get(ctx context.Context, name string, query string, dest any, args ...any) error {
	return ic.cmd.Get(ctx, name, query, dest, args...)
}
