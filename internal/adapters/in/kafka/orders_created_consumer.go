package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Haba1234/delivery/internal/core/application/usecases/commands"
	"github.com/Haba1234/delivery/internal/pkg/errs"
	"github.com/Haba1234/delivery/pkg/clients/basket/queues/basketconfirmedpb"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

type IOrdersCreateConsumer interface {
	Consume(ctx context.Context) error
	Close() error
}

var _ IOrdersCreateConsumer = &OrdersCreatedConsumer{}

type OrdersCreatedConsumer struct {
	topic                string
	consumer             *kafka.Consumer
	orderCreationHandler commands.ICreateOrderHandler
}

func NewOrdersCreateConsumer(
	host, group, topic string,
	orderCreationHandler commands.ICreateOrderHandler,
) (*OrdersCreatedConsumer, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}

	if group == "" {
		return nil, errs.NewValueIsRequiredError("group")
	}

	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}

	if orderCreationHandler == nil {
		return nil, errs.NewValueIsRequiredError("orderCreationHandler")
	}

	consumer, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":     host,
			"broker.address.family": "v4",
			"group.id":              group,
			"session.timeout.ms":    6000,
			"enable.auto.commit":    false,
			"auto.offset.reset":     "earliest",
			"client.id":             "delivery-service",
		},
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	return &OrdersCreatedConsumer{
		topic:                topic,
		consumer:             consumer,
		orderCreationHandler: orderCreationHandler,
	}, nil
}

func (c *OrdersCreatedConsumer) Close() error {
	return c.consumer.Close()
}

func (c *OrdersCreatedConsumer) Consume(ctx context.Context) error {
	err := c.consumer.Subscribe(c.topic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("context canceled")
			return nil
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Consumer error: %v (%v)\n", err, msg)
				continue
			}

			// Обрабатываем сообщение
			log.Printf("Received: %s => %s\n", msg.TopicPartition, string(msg.Value))
			var event basketconfirmedpb.BasketConfirmedIntegrationEvent

			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Отправляем команду
			newCreateOrder, err := commands.NewCreateOrder(uuid.New(), event.GetAddress().GetStreet())
			if err != nil {
				log.Printf("Failed to commands.NewCreateOrder: %v", err)
			}

			err = c.orderCreationHandler.Handle(ctx, newCreateOrder)
			if err != nil {
				log.Printf("Failed to orderCreationHandler.Handle: %v", err)
			}

			// Подтверждаем обработку сообщения
			_, err = c.consumer.CommitMessage(msg)
			if err != nil {
				log.Printf("Commit failed: %v", err)
			}
		}
	}
}
