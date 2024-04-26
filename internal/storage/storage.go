package storage

import (
	"TestWebServer/internal/model"
	"database/sql"
	"encoding/json"
	"log"
)

// IfDataExists checks if there is any data in the orders table.
// It returns true if there is data, false otherwise.
func (d *Database) IfDataExists() bool {
	var count int
	query := `SELECT COUNT(*) FROM orders`
	err := d.pg.Get(&count, query)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows were returned, which means the table is empty.
			log.Println("No data found in the orders table.")
			return false
		} else {
			// An error occurred while querying the database.
			log.Printf("Error checking for data: %v\n", err)
			return false
		}
	}
	// If count is greater than 0, there is data in the table.
	return count > 0
}

// IfStored returns true if order with a particular uid already exists in db and false otherwise
func (d *Database) IfStored(uid string) bool {
	var count int

	// Use a parameterized query to safely include the uid in the SQL statement
	query := `SELECT COUNT(*) FROM orders WHERE order_uid = $1`
	err := d.pg.Get(&count, query, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows were returned, which means the table is empty.
			log.Println("No data found for this uid")
			return false
		} else {
			// An error occurred while querying the database.
			log.Printf("Error checking for data: %v\n", err)
			return false
		}
	}
	// If count is greater than 0, there is data in the table.
	log.Printf("Data already stored")
	return count > 0
}

// FetchAllOrders retrieves all orders from the database, converts them to JSON, and returns a slice of JSON strings.
func (d *Database) FetchAllOrders() ([][]byte, error) {
	var orders []model.Order
	query := `SELECT * FROM orders`
	err := d.pg.Select(&orders, query)
	if err != nil {
		log.Printf("Error fetching orders: %v\n", err)
		return nil, err
	}

	var jsonOrders [][]byte
	for _, order := range orders {
		jsonOrder, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshaling order to JSON: %v\n", err)
			continue // Skip this order if it cannot be marshaled
		}
		jsonOrders = append(jsonOrders, jsonOrder)
	}

	return jsonOrders, nil
}

// InsertOrder inserts a new order into the database.
// It takes an Order model as input and executes an INSERT SQL query.
func (d *Database) InsertOrder(order model.Order) {
	_, err := d.pg.Exec(`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		*order.OrderUid, *order.TrackNumber, *order.Entry, *order.Locale, *order.InternalSignature, *order.CustomerId, *order.DeliveryService, *order.ShardKey, *order.SmId, *order.DateCreated, *order.OofShard)
	log.Print("ORDER INSERTED", *order.OrderUid)
	if err != nil {
		log.Print(err)
	}
}

// InsertDelivery inserts a new delivery record associated with an order into the database.
// It takes an order UID and a Delivery model as input and executes an INSERT SQL query.
func (d *Database) InsertDelivery(orderUid string, delivery model.Delivery) {
	_, err := d.pg.Exec(`INSERT INTO deliveries (order_id, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		orderUid, *delivery.Name, *delivery.Phone, *delivery.Zip, *delivery.City, *delivery.Address, *delivery.Region, *delivery.Email)
	if err != nil {
		log.Print(err)
	}
}

// InsertPayment inserts a new payment record associated with an order into the database.
// It takes an order UID and a Payment model as input and executes an INSERT SQL query.
func (d *Database) InsertPayment(orderUid string, payment model.Payment) {
	_, err := d.pg.Exec(`INSERT INTO payments (order_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		orderUid, *payment.Transaction, *payment.RequestId, *payment.Currency, *payment.Provider, *payment.Amount, *payment.PaymentDt, *payment.Bank, *payment.DeliveryCost, *payment.GoodsTotal, *payment.CustomFee)
	if err != nil {
		log.Print(err)
	}
}

// InsertItems inserts multiple item records associated with an order into the database.
// It takes an order UID and a slice of Item models as input and executes an INSERT SQL query for each item.
func (d *Database) InsertItems(orderUid string, items []model.Item) {
	for _, item := range items {
		_, err := d.pg.Exec(`INSERT INTO items (order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			orderUid, *item.ChrtId, *item.TrackNumber, *item.Price, *item.Rid, *item.Name, *item.Sale, *item.Size, *item.TotalPrice, *item.NmId, *item.Brand, *item.Status)
		if err != nil {
			log.Print(err)
		}
	}
}
