package broker

import "context"

type Consumer interface {
	Init() error
	ReadOrderMsg(ctx context.Context) error
	ReadRetryMsg(ctx context.Context) error
	ShutDown() error
	Ready() <-chan struct{}
}
