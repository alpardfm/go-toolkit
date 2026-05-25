package nosql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClose_CancelledContext(t *testing.T) {
	t.Run("returns error immediately when context is already cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		m := &mongoDB{}

		err := m.Close(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context already cancelled")
	})

	t.Run("returns error immediately when context deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		// Wait for deadline to pass
		<-ctx.Done()

		m := &mongoDB{}

		err := m.Close(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context already cancelled")
	})
}
