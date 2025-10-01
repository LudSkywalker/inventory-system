package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/LudSkywalker/inventory-system/store-service/domain/event"
	"github.com/Shopify/sarama"
)

const (
	TopicInventoryChanges = "inventory.changes"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("creating kafka producer: %w", err)
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) PublishInventoryChange(ctx context.Context, event event.InventoryEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: TopicInventoryChanges,
		Key:   sarama.StringEncoder(fmt.Sprintf("%s-%s", event.StoreID, event.ItemID)),
		Value: sarama.ByteEncoder(value),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
