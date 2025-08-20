package postgres

import (
	"context"
	"database/sql"
	"order-service/internal/domain"
)

func (p *PostgresDB) saveOrderTx(ctx context.Context, tx *sql.Tx, order *domain.Order) (int, error) {
	insertOrderQuery := `INSERT INTO orders 
	(order_uid, track_number, entry, customer_id, delivery_service, 
	date_created, locale, internal_signature, shardkey, sm_id, oof_shard)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	RETURNING id`
	var orderId int
	err := tx.QueryRowContext(ctx, insertOrderQuery,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.CustomerID,
		order.DeliveryService,
		order.DateCreated,
		order.Locale,
		order.InternalSignature,
		order.Shardkey,
		order.SmID,
		order.OofShard).Scan(&orderId)
	if err != nil {
		return -1, err
	}

	return orderId, nil
}

func (p *PostgresDB) saveDeliveryTx(ctx context.Context, tx *sql.Tx, orderId int, delivery *domain.Delivery) error {
	insertDeliveryQuery := `INSERT INTO deliveries
	(order_id, name, phone, zip, city, address, region, email)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := tx.ExecContext(ctx, insertDeliveryQuery,
		orderId,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresDB) savePaymentsTx(ctx context.Context, tx *sql.Tx, orderId int, payment *domain.Payment) error {
	insertPaymentQuery := `INSERT INTO payments
	(order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	_, err := tx.ExecContext(ctx, insertPaymentQuery,
		orderId,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)
	return err
}

func (p *PostgresDB) saveItemsTx(ctx context.Context, tx *sql.Tx, orderId int, items []*domain.Item) error {
	insertItemQuery := `INSERT INTO order_items
	(order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	stmt, err := tx.PrepareContext(ctx, insertItemQuery)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()
	for _, item := range items {
		_, err := stmt.ExecContext(ctx,
			orderId,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
