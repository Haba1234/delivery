package kafka

import (
	"encoding/json"

	"github.com/Haba1234/delivery/internal/core/domain/model/order"
	"github.com/Haba1234/delivery/internal/pkg/errs"
	"github.com/Haba1234/delivery/pkg/queues/orderstatuschangedpb"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type OrderProducer struct {
	topic    string
	producer *kafka.Producer
}

func NewOrderProducer(broker, topic string) (*OrderProducer, error) {
	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": broker,
		},
	)
	if err != nil {
		return nil, err
	}

	return &OrderProducer{
		topic:    topic,
		producer: producer,
	}, nil
}

func (p *OrderProducer) Close() error {
	p.producer.Close()

	return nil
}

func (p *OrderProducer) Produce(event order.CompletedDomainEvent) error {
	if event.IsEmpty() {
		return errs.NewValueIsRequiredError("event")
	}

	eventPayload, err := p.mapDomainEventToPbEvent(event)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	return p.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
			Key:            []byte(event.ID().String()),
			Value:          bytes,
		},
		p.producer.Events(),
	)
}

func (*OrderProducer) mapDomainEventToPbEvent(
	domainEvent order.CompletedDomainEvent,
) (*orderstatuschangedpb.OrderStatusChangedIntegrationEvent, error) {
	status, ok := orderstatuschangedpb.OrderStatus_value[domainEvent.OrderStatus()]
	if !ok {
		return nil, errs.NewValueIsInvalidError("OrderStatus")
	}

	integrationEvent := orderstatuschangedpb.OrderStatusChangedIntegrationEvent{
		OrderId:     domainEvent.OrderID().String(),
		OrderStatus: orderstatuschangedpb.OrderStatus(status),
	}

	return &integrationEvent, nil
}
