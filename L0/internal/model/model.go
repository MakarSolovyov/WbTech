package model

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Order struct {
	Order_uid          string    `json:"order_uid" fake:"{uuid}"`
	Track_number       string    `json:"track_number"`
	Entry              string    `json:"entry"`
	Delivery           Delivery  `json:"delivery" fake:"skip"`
	Payment            Payment   `json:"payment" fake:"skip"`
	Items              []Item    `json:"items" fake:"skip"`
	Locale             string    `json:"locale" fake:"{languageabbreviation}"`
	Internal_signature string    `json:"internal_signature" fake:"skip"`
	Customer_id        string    `json:"customer_id" fake:"{uuid}"`
	Delivery_service   string    `json:"delivery_service" fake:"{company}"`
	Shardkey           string    `json:"shardkey"`
	Sm_id              int       `json:"sm_id" fake:"{number:1,100}"`
	Date_created       time.Time `json:"date_created"`
	Oof_shard          string    `json:"oof_shard"`
}
type Delivery struct {
	Name    string `json:"name" fake:"{name}"`
	Phone   string `json:"phone" fake:"{phone}"`
	Zip     string `json:"zip"`
	City    string `json:"city" fake:"{city}"`
	Address string `json:"address" fake:"{address}"`
	Region  string `json:"region" fake:"{state}"`
	Email   string `json:"email" fake:"{email}"`
}
type Payment struct {
	Transaction   string  `json:"transaction" fake:"{uuid}"`
	Request_id    string  `json:"request_id" fake:"{uuid}"`
	Currency      string  `json:"currency" fake:"{currencylong}"`
	Provider      string  `json:"provider" fake:"{company}"`
	Amount        float64 `json:"amount" fake:"{price:1,1000000}"`
	Payment_dt    int     `json:"payment_dt" fake:"{number:100,100000}"`
	Bank          string  `json:"bank" fake:"{bankname}"`
	Delivery_cost float64 `json:"delivery_cost" fake:"{price:1,1000000}"`
	Goods_total   float64 `json:"goods_total" fake:"{price:1,1000000}"`
	Custom_fee    float64 `json:"custom_fee" fake:"{price:1,1000000}"`
}
type Item struct {
	Chrt_id      int     `json:"chrt_id" fake:"{number:1,10000000}"`
	Track_number string  `json:"track_number"`
	Price        float64 `json:"price" fake:"{price:1,1000}"`
	Rid          string  `json:"rid" fake:"{uuid}"`
	Name         string  `json:"name" fake:"{productname}"`
	Sale         int     `json:"sale" fake:"{number:0,100}"`
	Size         string  `json:"size"`
	Total_price  float64 `json:"total_price" fake:"{price:1,1000000}"`
	Nm_id        int     `json:"nm_id" fake:"{number:1,10000}"`
	Brand        string  `json:"brand" fake:"{company}"`
	Status       int     `json:"status" fake:"{number:0,10}"`
}

func ConnectDatabase() (*sql.DB, error) {

	dbinfo := "host=localhost port=5432 user=goAdmin dbname=dbgo password=password sslmode=disable"
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return nil, err
	}

	return db, nil
}

func AddToDatabase(order Order) error {

	db, err := ConnectDatabase()
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}

	err = AddOrder(tx, order)
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetOrderById(orderID string) (Order, error) {

	var order Order
	db, err := ConnectDatabase()
	if err != nil {
		log.Println(err)
		return order, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return order, err
	}

	order, err = GetOrder(tx, orderID)
	if err != nil {
		log.Println(err)
		return order, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return order, err
	}

	return order, nil
}

func GetOrders(orderCount int) ([]Order, error) {

	var orders []Order
	db, err := ConnectDatabase()
	if err != nil {
		log.Println(err)
		return orders, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return orders, err
	}

	rows, err := tx.Query(`Select * from order_service.orders
						   join order_service.delivery on order_service.delivery.delivery_id = order_service.orders.delivery 
						   join order_service.payment on order_service.payment.transaction = order_service.orders.payment
						   order by date_created desc limit $1`, orderCount)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return orders, err
	}
	defer rows.Close()

	var itemNmbers []string
	var ordersItems map[string]string = make(map[string]string)
	for rows.Next() {
		var order Order
		var delivery string
		var payment string
		var deliveryID string
		var itemsIDs string

		err := rows.Scan(&order.Order_uid, &order.Track_number, &order.Entry, &delivery, &payment,
			&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service,
			&order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard, &itemsIDs,

			&deliveryID, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
			&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,

			&order.Payment.Transaction, &order.Payment.Request_id, &order.Payment.Currency, &order.Payment.Provider,
			&order.Payment.Amount, &order.Payment.Payment_dt, &order.Payment.Bank, &order.Payment.Delivery_cost,
			&order.Payment.Goods_total, &order.Payment.Custom_fee)
		if err != nil {
			log.Println(err)
			return orders, err
		}

		itemNmbers = append(itemNmbers, itemsIDs)
		orders = append(orders, order)
		ordersItems[order.Order_uid] = itemsIDs
	}

	itemNmbersUnique := uniqueNumbers(itemNmbers)

	itemsMap, err := GetItemsMap(tx, itemNmbersUnique)
	if err != nil {
		log.Println(err)
		return orders, err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return orders, err
	}

	for index, value := range orders {
		var items []Item

		itemNumbers := strings.Split(ordersItems[value.Order_uid], ",")
		for _, num := range itemNumbers {
			intNumber, _ := strconv.Atoi(num)
			items = append(items, itemsMap[intNumber])
		}

		orders[index].Items = items
	}

	return orders, nil
}

func AddOrder(tx *sql.Tx, order Order) error {

	deliveryId, err := AddDelivery(tx, order.Delivery)
	if err != nil {
		log.Println(err)
		return err
	}

	paymentTransaction, err := AddPayment(tx, order.Payment)
	if err != nil {
		log.Println(err)
		return err
	}

	var itemIdArray []int32
	for _, value := range order.Items {
		itemId, err := AddItem(tx, value)
		if err != nil {
			log.Println(err)
			return err
		}
		itemIdArray = append(itemIdArray, itemId)
	}

	_, err = tx.Exec(`Insert into order_service.orders 
					  (order_uid, track_number, entry, delivery, payment,
					   items, locale, internal_signature, customer_id,
					   delivery_service, shardkey, sm_id, date_created, oof_shard) 
					   values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		order.Order_uid, order.Track_number,
		order.Entry, deliveryId,
		paymentTransaction, int32SliceToString(itemIdArray),
		order.Locale, order.Internal_signature,
		order.Customer_id, order.Delivery_service,
		order.Shardkey, order.Sm_id,
		order.Date_created, order.Oof_shard)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return err
	}

	log.Printf("Заказ (uid: %#v) добавлен в БД.", order.Order_uid)
	return nil
}

func AddDelivery(tx *sql.Tx, delivery Delivery) (string, error) {
	deliveryId := uuid.New().String()

	_, err := tx.Exec(`Insert into order_service.delivery 
					  (delivery_id, name, phone, zip, city, address, region, email) 
					   values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		deliveryId, delivery.Name,
		delivery.Phone, delivery.Zip,
		delivery.City, delivery.Address,
		delivery.Region, delivery.Email)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return "", nil
	}

	return deliveryId, nil
}

func AddPayment(tx *sql.Tx, payment Payment) (string, error) {

	_, err := tx.Exec(`Insert into order_service.payment 
					  (transaction, request_id, currency, provider, amount,
					   payment_dt, bank, delivery_cost, goods_total, custom_fee) 
					   values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		payment.Transaction, payment.Request_id,
		payment.Currency, payment.Provider,
		payment.Amount, payment.Payment_dt,
		payment.Bank, payment.Delivery_cost,
		payment.Goods_total, payment.Custom_fee)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return "", nil
	}

	return payment.Transaction, nil
}

func AddItem(tx *sql.Tx, item Item) (int32, error) {

	_, err := tx.Exec(`Insert into order_service.items 
					  (chrt_id, track_number, price, rid, name,
					   sale, size, total_price, nm_id, brand, status) 
					   values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		item.Chrt_id, item.Track_number,
		item.Price, item.Rid,
		item.Name, item.Sale,
		item.Size, item.Total_price,
		item.Nm_id, item.Brand,
		item.Status)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return 0, nil
	}

	return int32(item.Chrt_id), nil
}

func GetOrder(tx *sql.Tx, orderID string) (Order, error) {

	var order Order
	var delivery string
	var payment string
	var deliveryID string
	var itemsIDs string

	row := tx.QueryRow(`Select * from order_service.orders 
	join order_service.delivery on order_service.delivery.delivery_id = order_service.orders.delivery 
	join order_service.payment on order_service.payment.transaction = order_service.orders.payment
	where order_service.orders.order_uid = $1`, orderID)

	err := row.Scan(&order.Order_uid, &order.Track_number, &order.Entry, &delivery, &payment,
		&order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service,
		&order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard, &itemsIDs,

		&deliveryID, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,

		&order.Payment.Transaction, &order.Payment.Request_id, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &order.Payment.Payment_dt, &order.Payment.Bank, &order.Payment.Delivery_cost,
		&order.Payment.Goods_total, &order.Payment.Custom_fee)
	if err != nil {
		log.Println(err)
		return order, err
	}

	orderitems, err := GetItems(tx, itemsIDs)
	if err != nil {
		log.Println(err)
		return order, err
	}

	order.Items = orderitems

	return order, nil
}

// TODO: GetItems и GetItemsMap - объединить в одну функцию
func GetItems(tx *sql.Tx, itemsIDs string) ([]Item, error) {
	var itemArray []Item

	chrtIDs, _ := stringToInt32Slice(itemsIDs)
	rows, err := tx.Query(`SELECT * FROM order_service.items WHERE chrt_id = any($1)`, pq.Array(chrtIDs))
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return itemArray, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Item

		err := rows.Scan(&item.Chrt_id, &item.Track_number, &item.Price, &item.Rid,
			&item.Name, &item.Sale, &item.Size, &item.Total_price,
			&item.Nm_id, &item.Brand, &item.Status)
		if err != nil {
			log.Println(err)
			return itemArray, err
		}

		itemArray = append(itemArray, item)
	}

	return itemArray, nil
}

func GetItemsMap(tx *sql.Tx, itemsIDs string) (map[int]Item, error) {
	var itemsMap map[int]Item = make(map[int]Item)

	chrtIDs, _ := stringToInt32Slice(itemsIDs)
	rows, err := tx.Query(`SELECT * FROM order_service.items WHERE chrt_id = any($1)`, pq.Array(chrtIDs))
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return itemsMap, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Item

		err := rows.Scan(&item.Chrt_id, &item.Track_number, &item.Price, &item.Rid,
			&item.Name, &item.Sale, &item.Size, &item.Total_price,
			&item.Nm_id, &item.Brand, &item.Status)
		if err != nil {
			log.Println(err)
			return itemsMap, err
		}

		itemsMap[item.Chrt_id] = item
	}

	return itemsMap, nil
}

func int32SliceToString(nums []int32) string {
	var strNums []string
	for _, num := range nums {
		strNums = append(strNums, fmt.Sprint(num))
	}
	return strings.Join(strNums, ",")
}

func stringToInt32Slice(s string) ([]int32, error) {
	ids := strings.Split(s, ",")
	result := make([]int32, len(ids))
	for i, id := range ids {
		n, err := strconv.Atoi(strings.TrimSpace(id))
		if err != nil {
			return nil, err
		}
		result[i] = int32(n)
	}
	return result, nil
}

func uniqueNumbers(strs []string) string {
	unique := make(map[string]bool)
	for _, str := range strs {
		numbers := strings.Split(str, ",")
		for _, num := range numbers {
			trimmedNum := strings.TrimSpace(num)
			if trimmedNum != "" {
				unique[trimmedNum] = true
			}
		}
	}

	var result []string
	for key := range unique {
		result = append(result, key)
	}
	return strings.Join(result, ",")
}
