# Graceful Shutdown Guide

This guide documents how to gracefully shut down database connections using the `go-toolkit` SQL and NoSQL packages. A proper shutdown ensures in-flight operations complete cleanly before the application exits.

## Method Signatures

### SQL Package

```go
import "github.com/alpardfm/go-toolkit/heavy/sql"

type Interface interface {
    Leader() Command
    Follower() Command
    Stop()
    Shutdown(ctx context.Context) error
}
```

**`Shutdown(ctx context.Context) error`**

Gracefully closes database connections, waiting for in-flight queries to complete until the context deadline is reached.

- If no deadline is set on the provided context, a default timeout of **30 seconds** is applied.
- Returns `nil` if all in-flight operations complete before the deadline.
- Returns a timeout error (code `CodeSQLShutdown`) if the deadline is reached while operations are still active.
- After shutdown is initiated, no new queries are accepted.

### NoSQL Package

```go
import "github.com/alpardfm/go-toolkit/heavy/nosql"

type Interface interface {
    Close(ctx context.Context) error
    // ...other methods
}
```

**`Close(ctx context.Context) error`**

Disconnects from MongoDB, respecting the context deadline.

- If the context is already cancelled, returns an error (code `CodeNoSQLClose`) without attempting disconnection.
- Propagates the context deadline to the underlying MongoDB `Disconnect` call, giving in-flight operations time to complete.
- Returns `nil` on successful disconnection.

## Registering OS Signal Handlers

Use Go's `os/signal` package to listen for termination signals (`SIGTERM`, `SIGINT`) and trigger graceful shutdown:

```go
import (
    "context"
    "os"
    "os/signal"
    "syscall"
)

// Create a channel to receive OS signals
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

// Block until a signal is received
<-quit
```

When deploying in containers or orchestrated environments (Kubernetes, ECS), `SIGTERM` is the standard signal sent before forced termination. Registering both `SIGTERM` and `SIGINT` covers production deployments and local development (Ctrl+C).

## Shutdown Order of Operations

Follow this sequence to ensure a clean shutdown:

1. **Stop accepting new requests** — Close the HTTP server or stop consuming from message queues.
2. **Wait for in-flight requests to drain** — Allow currently processing requests to finish (HTTP server's `Shutdown` handles this).
3. **Close database connections** — Shut down SQL and NoSQL connections with a deadline context.
4. **Exit the process** — Return from `main()` or call `os.Exit(0)`.

```
Signal received
    │
    ▼
┌─────────────────────────────┐
│ 1. Stop accepting requests  │  ← http.Server.Shutdown(ctx)
└─────────────────────────────┘
    │
    ▼
┌─────────────────────────────┐
│ 2. Wait for in-flight       │  ← Handled by HTTP server shutdown
│    requests to complete     │
└─────────────────────────────┘
    │
    ▼
┌─────────────────────────────┐
│ 3. Close SQL connections    │  ← sqlDB.Shutdown(ctx)
└─────────────────────────────┘
    │
    ▼
┌─────────────────────────────┐
│ 4. Close NoSQL connections  │  ← mongoDB.Close(ctx)
└─────────────────────────────┘
    │
    ▼
┌─────────────────────────────┐
│ 5. Exit                     │
└─────────────────────────────┘
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/alpardfm/go-toolkit/heavy/log"
    "github.com/alpardfm/go-toolkit/heavy/nosql"
    "github.com/alpardfm/go-toolkit/heavy/sql"
)

func main() {
    // Initialize logger
    logger := log.Init(log.Config{Level: "info"})

    // Initialize SQL database
    sqlDB, err := sql.Init(sql.Config{
        Driver: "postgres",
        Leader: sql.ConnConfig{
            Host:     "localhost",
            Port:     5432,
            DB:       "myapp",
            User:     "user",
            Password: "pass",
        },
    }, logger)
    if err != nil {
        logger.Fatal(context.Background(), fmt.Sprintf("failed to init SQL: %v", err))
    }

    // Initialize NoSQL database
    noSQL := nosql.Init(nosql.Config{
        DBUrl: "mongodb://localhost:27017",
        DB:    "myapp",
    }, logger)

    // Set up HTTP server
    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Start server in a goroutine
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal(context.Background(), fmt.Sprintf("server error: %v", err))
        }
    }()

    logger.Info(context.Background(), "server started on :8080")

    // Wait for termination signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    logger.Info(context.Background(), "shutting down...")

    // Create a deadline context for the entire shutdown sequence
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Step 1 & 2: Stop accepting new requests and wait for in-flight to drain
    if err := server.Shutdown(ctx); err != nil {
        logger.Error(ctx, fmt.Errorf("HTTP server shutdown error: %w", err))
    }

    // Step 3: Close SQL connections gracefully
    if err := sqlDB.Shutdown(ctx); err != nil {
        logger.Error(ctx, fmt.Errorf("SQL shutdown error: %w", err))
    }

    // Step 4: Close NoSQL connections
    if err := noSQL.Close(ctx); err != nil {
        logger.Error(ctx, fmt.Errorf("NoSQL close error: %w", err))
    }

    logger.Info(context.Background(), "shutdown complete")
}
```

## Tips

- **Set appropriate timeouts.** The SQL `Shutdown` defaults to 30 seconds if no deadline is set, but you should set an explicit deadline that fits your deployment environment (e.g., Kubernetes gives 30s by default before `SIGKILL`).
- **Shut down in dependency order.** Close the HTTP server first so no new requests arrive that need database access, then close databases.
- **Log shutdown progress.** Logging each step helps diagnose slow shutdowns in production.
- **Handle errors but don't block.** If one connection fails to close, log the error and continue shutting down other resources.
