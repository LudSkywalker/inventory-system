package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

	// Configure SASL authentication
	if saslEnable := os.Getenv("SASL_ENABLE"); saslEnable == "true" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = os.Getenv("SASL_USER")
		config.Net.SASL.Password = os.Getenv("SASL_PASSWORD")

		if mechanism := os.Getenv("SASL_MECHANISM"); mechanism != "" {
			config.Net.SASL.Mechanism = sarama.SASLMechanism(mechanism)
		} else {
			config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		}

		// Configure TLS for SASL
		if tlsEnable := os.Getenv("TLS_ENABLE"); tlsEnable == "true" {
			config.Net.TLS.Enable = true
		}
	}

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
