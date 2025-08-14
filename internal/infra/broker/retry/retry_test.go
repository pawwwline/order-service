package retry

import (
	"context"
	"order-service/internal/config"
	"order-service/internal/infra/broker/handler"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetryWrapperResults(t *testing.T) {
	retry := NewRetry(config.KafkaConfig{
		RetryMaxAttempts:   5,
		BackoffDurationMin: 10,
		BackoffDurationMax: 100,
	})
	ctx := context.Background()

	t.Run("always retryable", func(t *testing.T) {
		attempts := 0
		res := retry.RetryWrapper(ctx, func() handler.Result {
			attempts++
			return handler.Retry
		})
		assert.Equal(t, handler.Retry, res)
		assert.Equal(t, retry.MaxAttempts(), attempts)
	})

	t.Run("success first try", func(t *testing.T) {
		attempts := 0
		res := retry.RetryWrapper(ctx, func() handler.Result {
			attempts++
			return handler.Success
		})
		assert.Equal(t, handler.Success, res)
		assert.Equal(t, 1, attempts)
	})

	t.Run("success after retries", func(t *testing.T) {
		attempts := 0
		res := retry.RetryWrapper(ctx, func() handler.Result {
			attempts++
			if attempts < retry.MaxAttempts() {
				return handler.Retry
			}
			return handler.Success
		})
		assert.Equal(t, handler.Success, res)
		assert.Equal(t, retry.MaxAttempts(), attempts)
	})

	t.Run("DLQ returned immediately", func(t *testing.T) {
		attempts := 0
		res := retry.RetryWrapper(ctx, func() handler.Result {
			attempts++
			return handler.DLQ
		})
		assert.Equal(t, handler.DLQ, res)
		assert.Equal(t, 1, attempts)
	})
}
