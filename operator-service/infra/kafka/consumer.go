package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
			var event event.InventoryEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				continue
			}

			dto := dto.GlobalInventoryDTO{
				ItemID:    event.ItemID,
				StoreID:   event.StoreID,
				Quantity:  event.Quantity,
				UpdatedAt: event.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			}

			if err := c.useCase.ProcessInventoryEvent(ctx, dto); err != nil {
				log.Printf("Error processing event: %v", err)
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
