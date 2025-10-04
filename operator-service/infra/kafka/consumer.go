package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/LudSkywalker/inventory-system/operator-service/app/dto"
	"github.com/LudSkywalker/inventory-system/operator-service/app/port/input"
	"github.com/LudSkywalker/inventory-system/operator-service/domain/event"
	"github.com/Shopify/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	useCase  input.GlobalInventoryUseCase
}

func NewConsumer(brokers []string, useCase input.GlobalInventoryUseCase) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

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

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("creating kafka consumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		useCase:  useCase,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	partitionConsumer, err := c.consumer.ConsumePartition("inventory.changes", 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("starting partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	log.Println("Kafka consumer started")

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			if err := c.processMessage(ctx, msg); err != nil {
				log.Printf("Error processing message: %v", err)
			}

		case err := <-partitionConsumer.Errors():
			log.Printf("Error from consumer: %v", err)

		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func (c *Consumer) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	// Validate message payload
	if msg.Value == nil || len(msg.Value) == 0 {
		return fmt.Errorf("received empty or nil message at partition %d, offset %d", msg.Partition, msg.Offset)
	}

	var event event.InventoryEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("Failed to unmarshal event, raw JSON data: %s", string(msg.Value))
		return fmt.Errorf("failed to unmarshal event at partition %d, offset %d: %w", msg.Partition, msg.Offset, err)
	}

	// Basic validation of the event
	if event.ItemID == "" {
		return fmt.Errorf("invalid event: missing ItemID at partition %d, offset %d", msg.Partition, msg.Offset)
	}

	// Convert to DTO
	dto := dto.GlobalInventoryDTO{
		ItemID:    event.ItemID,
		ItemName:  event.ItemName,
		StoreID:   event.StoreID,
		Quantity:  event.Quantity,
		UpdatedAt: event.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Process the event
	return c.useCase.ProcessInventoryEvent(ctx, dto)
}
