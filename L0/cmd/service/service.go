package main

import (
	"l0/internal/model"
	kafka "l0/internal/pkg"
	validation "l0/internal/validator"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

func main() {

	log.Println("Сервис запущен.")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	//var brokers []string = []string{"localhost:9092"}
	brokers := []string{os.Getenv("KAFKA_BROKERS")} // "kafka:9092"

	consumer, err := kafka.CreateConsumer(brokers)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	producer, err := kafka.CreateProducer(brokers)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	doneCh := make(chan struct{})

	go kafka.DoServiceRequest(producer, consumer, doneCh, generateOrdersOperation, "post_order", "post_order_response")
	go kafka.DoServiceRequest(producer, consumer, doneCh, getOrderByIdOperation, "get_order_by_id", "get_order_by_id_response")
	go stopService(sigchan, doneCh)

	<-doneCh
}

func generateOrdersOperation(data string) (interface{}, error) {

	orderCount, err := strconv.Atoi(data)
	if err != nil {
		log.Printf("Ошибка преобразования: %v\n", err)
		return nil, err
	}

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	var ordersAdded int = 0
	for i := 0; i < orderCount; i++ {

		var order model.Order
		var delivery model.Delivery
		var payment model.Payment
		var items []model.Item

		err := gofakeit.Struct(&order)
		if err != nil {
			log.Printf("Ошибка генерации: %v\n", err)
			return nil, err
		}

		err = gofakeit.Struct(&delivery)
		if err != nil {
			log.Printf("Ошибка генерации: %v\n", err)
			return nil, err
		}

		err = gofakeit.Struct(&payment)
		if err != nil {
			log.Printf("Ошибка генерации: %v\n", err)
			return nil, err
		}

		itemsNumber := rng.Intn(10) + 1
		for i := 0; i < itemsNumber; i++ {

			var item model.Item
			err = gofakeit.Struct(&item)
			if err != nil {
				log.Printf("Ошибка генерации: %v\n", err)
				return nil, err
			}

			items = append(items, item)
		}

		order.Delivery = delivery
		order.Payment = payment
		order.Items = items

		if err = validation.Val.Struct(&order); err != nil {
			log.Printf("Заказ (uuid: %s) не прошел валидацию: %s", order.Order_uid, err)
			return nil, err
		}

		// TODO: Нужно сделать запись нескольких заказов одним запросом
		err = model.AddToDatabase(order)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		ordersAdded++
	}

	message := "Сгенерировано заказов: " + strconv.Itoa(ordersAdded)
	return message, nil
}

func getOrderByIdOperation(uid string) (interface{}, error) {

	uid = uid[1 : len(uid)-1]

	order, err := model.GetOrderById(uid)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return order, nil
}

func stopService(sigchan chan os.Signal, doneCh chan struct{}) {
	<-sigchan
	log.Println("Горячая клавиша выхода зафиксирована, сервис остановлен.")
	doneCh <- struct{}{}
}
