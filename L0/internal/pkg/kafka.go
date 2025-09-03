package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type InitController struct {
	Producer sarama.SyncProducer
	Worker   sarama.Consumer
}

func CreateProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second

	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Errors = true

	config.Producer.Timeout = 5 * time.Second

	return sarama.NewSyncProducer(brokers, config)
}

func CreateConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}

func DoRequest[T any](producer sarama.SyncProducer, worker sarama.Consumer, object T, topic string, topicRespond string) string {

	err := SendMessage(producer, object, topic)
	if err != nil {
		//log.Println(err)
		return err.Error()
	}

	consumer, err := worker.ConsumePartition(topicRespond, 0, sarama.OffsetNewest)
	if err != nil {
		//log.Println(err)
		return err.Error()
	}

	timer := time.After(10 * time.Second)

	var response string = ""
	doneCh := make(chan struct{})
	go func() {
		select {
		case err := <-consumer.Errors():

			var consError error = err.Err
			//log.Println(consError)
			response = consError.Error()
			consumer.Close()

			doneCh <- struct{}{}
		case msg := <-consumer.Messages():

			log.Printf("Message received: Topic(%s) | Message(%s) \n",
				string(msg.Topic),
				string(msg.Value))

			response = string(msg.Value)
			consumer.Close()

			doneCh <- struct{}{}
		case <-timer:

			response = "Таймаут ожидания сообщения истек"
			log.Println(response)
			consumer.Close()

			doneCh <- struct{}{}
		}
	}()
	<-doneCh

	return response
}

func DoServiceRequest(producer sarama.SyncProducer, worker sarama.Consumer, doneCh chan struct{},
	operation func(string) (interface{}, error), topic string, topicRespond string) {

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():

				var consError error = err.Err
				log.Println(consError)
				if consError != nil {
					result := consError.Error()

					errMsg := SendMessage(producer, result, topicRespond)
					if err != nil {
						log.Println(errMsg)
						continue
					}
				}

			case msg := <-consumer.Messages():

				log.Printf("Message received: Topic(%s) | Message(%s) \n",
					string(msg.Topic),
					string(msg.Value))

				result, err := operation(string(msg.Value))
				if err != nil {
					log.Println(err)
					result = err.Error()
				}

				err = SendMessage(producer, result, topicRespond)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()

	<-doneCh
}

func SendMessage[T any](producer sarama.SyncProducer, order T, topic string) error {

	msgInBytes, err := json.Marshal(order)
	if err != nil {
		log.Println(err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msgInBytes),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Message sent: Topic(%s) | Partition(%d) | offset(%d)\n",
		topic,
		partition,
		offset)

	return nil
}
