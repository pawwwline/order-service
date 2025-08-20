package postgres

import (
	"context"
	"order-service/internal/domain"
)

func (p *PostgresDB) getOrderItems(ctx context.Context, orderId int) ([]*domain.Item, error) {
	var items []*domain.Item
	rows, err := p.db.QueryContext(ctx, `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM order_items WHERE order_id = $1`, orderId)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var item domain.Item
		if err := rows.Scan(&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil

}

func (p *PostgresDB) getItemsByOrderIds(ctx context.Context, orderIdArr []int) (map[int][]*domain.Item, error) {
	query := `
    SELECT order_id, chrt_id, track_number, price, rid, name, sale, size,
           total_price, nm_id, brand, status
    FROM order_items
    WHERE order_id = ANY($1)
`
	rows, err := p.db.QueryContext(ctx, query, orderIdArr)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	itemsByOrder := make(map[int][]*domain.Item)

	for rows.Next() {
		var (
			orderID int
			item    domain.Item
		)
		if err := rows.Scan(
			&orderID, &item.ChrtID, &item.TrackNumber, &item.Price,
			&item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		); err != nil {
			return nil, err
		}
		itemsByOrder[orderID] = append(itemsByOrder[orderID], &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return itemsByOrder, nil
}

func (p *PostgresDB) getOrdersIds(orders []*domain.Order) []int {
	ids := make([]int, 0, len(orders))
	for _, o := range orders {
		ids = append(ids, o.Id)
	}
	return ids
}

func (p *PostgresDB) getOrdersWithoutItems(ctx context.Context, limit int) ([]*domain.Order, error) {
	rows, err := p.db.QueryContext(ctx, `
		SELECT 
			o.id, o.order_uid, o.track_number, o.entry, o.customer_id, o.delivery_service,
			o.date_created, o.locale, o.internal_signature, o.shardkey, o.sm_id, o.oof_shard,
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
			p.transaction, p.request_id, p.provider, p.amount, p.payment_dt, p.bank,
			p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		JOIN deliveries d ON o.id = d.order_id
		JOIN payments   p ON o.id = p.order_id
		ORDER BY o.date_created DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		var delivery domain.Delivery
		var payment domain.Payment

		if err := rows.Scan(
			&o.Id, &o.OrderUID, &o.TrackNumber, &o.Entry, &o.CustomerID, &o.DeliveryService,
			&o.DateCreated, &o.Locale, &o.InternalSignature, &o.Shardkey, &o.SmID, &o.OofShard,
			&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
			&payment.Transaction, &payment.RequestID, &payment.Provider, &payment.Amount, &payment.PaymentDt,
			&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
		); err != nil {
			return nil, err
		}

		o.Delivery = &delivery
		o.Payment = &payment
		orders = append(orders, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil

}

func (p *PostgresDB) attachItemsToOrder(orders []*domain.Order, items map[int][]*domain.Item) {
	for _, o := range orders {
		if items, ok := items[o.Id]; ok {
			o.Items = items
		}
	}
}
