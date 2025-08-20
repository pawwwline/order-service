package broker

import (
	"context"
	"sync"
)

type Broker struct {
	consumer Consumer
	wg       sync.WaitGroup
}

func NewBroker(consumer Consumer) *Broker {
	return &Broker{consumer: consumer}
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
			err := b.consumer.ReadOrderMsg(ctx)
			if err != nil {
				return
			}
		}
	}
}

func (b *Broker) runRetries(ctx context.Context) {
	defer b.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := b.consumer.ReadRetryMsg(ctx)
			if err != nil {
				return
			}
		}
	}
}
