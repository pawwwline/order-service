package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"order-service/internal/domain"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(db *sql.DB) *PostgresDB {
	return &PostgresDB{
		db: db,
	}
}

func (p *PostgresDB) SaveOrder(ctx context.Context, order *domain.Order) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	orderID, err := p.saveOrderTx(ctx, tx, order)
	if err != nil {
		tx.Rollback()
		return err
	}

	if order.Delivery != nil {
		if err := p.saveDeliveryTx(ctx, tx, orderID, order.Delivery); err != nil {
			tx.Rollback()
			return err
		}
	}

	if order.Payment != nil {
		if err := p.savePaymentsTx(ctx, tx, orderID, order.Payment); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(order.Items) > 0 {
		if err := p.saveItemsTx(ctx, tx, orderID, order.Items); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (p *PostgresDB) GetOrderByUid(ctx context.Context, orderUID string) (*domain.Order, error) {
	row := p.db.QueryRowContext(ctx, `
	SELECT 
    o.id, o.order_uid, o.track_number, o.entry, o.customer_id, o.delivery_service,
    o.date_created, o.date_updated, o.locale, o.internal_signature, o.shardkey, o.sm_id, o.oof_shard,
    d.id AS delivery_id, d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
    p.id AS payment_id, p.transaction, p.request_id, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
	FROM orders o
	JOIN deliveries d ON o.id = d.order_id
	JOIN payments p ON o.id = p.order_id
	WHERE o.order_uid = $1;
`, orderUID)
	var orderId int
	var order domain.Order
	var delivery domain.Delivery
	var payment domain.Payment

	err := row.Scan(
		&orderId, &order.OrderUID, &order.TrackNumber, &order.Entry, &order.CustomerID, &order.DeliveryService,
		&order.DateCreated, &order.Locale, &order.InternalSignature, &order.Shardkey, &order.SmID, &order.OofShard,
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
		&payment.Transaction, &payment.RequestID, &payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}

	order.Delivery = &delivery
	order.Payment = &payment

	items, err := p.getOrderItems(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("error getting items %w", err)
	}
	order.Items = items

	return &order, nil

}

func (p *PostgresDB) CheckIdempotencyKey(ctx context.Context, key string) (bool, error) {
	var exists bool
	err := p.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)`, key).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *PostgresDB) GetLastOrders(ctx context.Context, limit int) ([]*domain.Order, error) {
	orders, err := p.getOrdersWithoutItems(ctx, limit)
	if err != nil {
		return nil, err
	}

	orderIds := p.getOrdersIds(orders)

	orderItems, err := p.getItemsByOrderIds(ctx, orderIds)
	if err != nil {
		return nil, err
	}
	p.attachItemsToOrder(orders, orderItems)
	return orders, nil

}
