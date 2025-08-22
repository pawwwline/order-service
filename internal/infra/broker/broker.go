package broker

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"
)

type Broker struct {
	consumer Consumer
	wg       sync.WaitGroup
	logger   *slog.Logger
}

func NewBroker(consumer Consumer, logger *slog.Logger) *Broker {
	return &Broker{consumer: consumer, logger: logger}
}

func (b *Broker) Run(ctx context.Context) {
	if err := b.consumer.Init(); err != nil {
		return
	}

	<-b.consumer.Ready()

	b.wg.Add(2)
	go b.runOrders(ctx)
	go b.runRetries(ctx)

}

func (b *Broker) Shutdown() error {
	b.wg.Wait()

	return b.consumer.ShutDown()
}

func (b *Broker) runOrders(ctx context.Context) {
	defer b.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := b.consumer.ReadOrderMsg(ctx)
		if err == nil {
			continue
		}

		if b.retryableErr(err) {
			select {
			case <-time.After(time.Second):
			case <-ctx.Done():
				return
			}
			b.logger.Error("retry error", "err", err)

			continue
		}

		b.logger.Error("not kafka temporary error", "err", err)
	}
}

func (b *Broker) runRetries(ctx context.Context) {
	defer b.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := b.consumer.ReadRetryMsg(ctx)
		if err == nil {
			continue
		}

		if b.retryableErr(err) {
			select {
			case <-time.After(time.Second):
			case <-ctx.Done():
				return
			}
			b.logger.Error("retry error", "err", err)

			continue
		}

		b.logger.Error("not kafka temporary error", "err", err)
	}

}

func (b *Broker) retryableErr(err error) bool {
	var opErr *net.OpError

	if errors.As(err, &opErr) {
		msg := opErr.Err.Error()
		if strings.Contains(msg, "failed to open connection") ||
			strings.Contains(msg, "connection refused") ||
			strings.Contains(msg, "connection reset") {
			return true
		}
	}

	var dnsErr *net.DNSError

	if errors.As(err, &dnsErr) {
		return dnsErr.IsTimeout
	}

	return false
}
