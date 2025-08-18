package retry

import (
	"context"
	"math/rand"
	"order-service/internal/config"
	"order-service/internal/infra/broker/handler"
	"time"
)

type Retry struct {
	maxAttempts        int
	backoffDurationMin int
	backoffDurationMax int
}

func NewRetry(cfg *config.KafkaConfig) *Retry {
	return &Retry{
		maxAttempts:        cfg.RetryMaxAttempts,
		backoffDurationMin: cfg.BackoffDurationMin,
		backoffDurationMax: cfg.BackoffDurationMax,
	}
}

func (r *Retry) MaxAttempts() int {
	return r.maxAttempts
}

func (r *Retry) BackoffDurationMin() time.Duration {
	return time.Duration(r.backoffDurationMin)
}

// exponential backoff logic
func (r *Retry) BackoffDuration(attempt int) time.Duration {
	if attempt < 1 {
		return 0
	}

	backoff := r.backoffDurationMin * (1 << (attempt - 1))
	if backoff > r.backoffDurationMax {
		backoff = r.backoffDurationMax
	}

	// adding jitter to avoid thundering herd problem
	jitter := rand.Int63n(int64(backoff)/2) - int64(backoff)/4

	return time.Duration(int64(backoff)+jitter) * time.Millisecond
}

// decorator for retry logic
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
