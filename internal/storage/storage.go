package storage

import (
	"TestWebServer/internal/model"

	"log"
)

func (d *Database) InsertOrder(order model.Order) {
	_, err := d.pg.Exec(`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		*order.OrderUid, *order.TrackNumber, *order.Entry, *order.Locale, *order.InternalSignature, *order.CustomerId, *order.DeliveryService, *order.ShardKey, *order.SmId, *order.DateCreated, *order.OofShard)
	if err != nil {
		log.Print(err)
	}
}

func (d *Database) InsertDelivery(orderUid string, delivery model.Delivery) {
	_, err := d.pg.Exec(`INSERT INTO deliveries (order_id, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		orderUid, *delivery.Name, *delivery.Phone, *delivery.Zip, *delivery.City, *delivery.Address, *delivery.Region, *delivery.Email)
	if err != nil {
		log.Print(err)
	}
}

func (d *Database) InsertPayment(orderUid string, payment model.Payment) {
	_, err := d.pg.Exec(`INSERT INTO payments (order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		orderUid, *payment.Transaction, *payment.RequestId, *payment.Currency, *payment.Provider, *payment.Amount, *payment.PaymentDt, *payment.Bank, *payment.DeliveryCost, *payment.GoodsTotal, *payment.CustomFee)
	if err != nil {
		log.Print(err)
	}
}

func (d *Database) InsertItems(orderUid string, items []model.Item) {
	for _, item := range items {
		_, err := d.pg.Exec(`INSERT INTO items (order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			orderUid, *item.ChrtId, *item.TrackNumber, *item.Price, *item.Rid, *item.Name, *item.Sale, *item.Size, *item.TotalPrice, *item.NmId, *item.Brand, *item.Status)
		if err != nil {
			log.Print(err)
		}
	}
}
