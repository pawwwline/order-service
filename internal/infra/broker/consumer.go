package broker

import "context"

type Consumer interface {
	ReadOrderMsg(ctx context.Context)
	ReadRetryMsg(ctx context.Context)
	ShutDown() error
}
