package retry

import (
	"context"
	"math/rand"
	"order-service/internal/config"
	"order-service/internal/infra/broker/handler"
	"time"
)

type Retry struct {
	maxAttempts int
	backoffMin  time.Duration
	backoffMax  time.Duration
}

func NewRetry(cfg config.KafkaConfig) *Retry {
	return &Retry{
		maxAttempts: cfg.RetryMaxAttempts,
		backoffMin:  time.Duration(cfg.BackoffDurationMin) * time.Second,
		backoffMax:  time.Duration(cfg.BackoffDurationMax) * time.Second,
	}
}

func (r *Retry) BackoffDuration(attempt int) time.Duration {
	if attempt < 1 {
		return 0
	}

	backoff := r.backoffMin * (1 << (attempt - 1))
	if backoff > r.backoffMax {
		backoff = r.backoffMax
	}

	jitter := time.Duration(rand.Int63n(int64(backoff)/2) - int64(backoff)/4)
	return backoff + jitter
}

// decorator for retry logic.

func (r *Retry) RetryWrapper(ctx context.Context, fn func() handler.Result) handler.Result {
	var res handler.Result
	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return res
		default:

			res = fn()

			if res == handler.Success || res == handler.DLQ {
				return res
			}

			if attempt < r.maxAttempts {
				select {
				case <-time.After(r.BackoffDuration(attempt)):
				case <-ctx.Done():
					return res
				}
			}
		}
	}

	return res
}
