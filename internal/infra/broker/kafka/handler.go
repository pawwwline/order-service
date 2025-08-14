package kafka

import (
	"context"
	"order-service/internal/infra/broker/handler"
)

type Handler interface {
	ProcessOrderMessage(ctx context.Context, msg []byte) handler.Result
}

type RetryHandler interface {
	RetryWrapper(ctx context.Context, fn func() handler.Result) handler.Result
}
