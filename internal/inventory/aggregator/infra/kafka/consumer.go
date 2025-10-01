package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/LudSkywalker/inventory-system/internal/inventory/aggregator/app/service"
	"github.com/LudSkywalker/inventory-system/internal/inventory/core/event"
	"github.com/Shopify/sarama"
)

type Consumer struct {
	consumer  sarama.Consumer
	processor *service.InventoryService
	topicName string
}

func NewConsumer(brokers []string, processor *service.InventoryService) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("creating kafka consumer: %w", err)
	}

	return &Consumer{
		consumer:  consumer,
		processor: processor,
		topicName: "inventory.changes",
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topicName, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("starting partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	log.Printf("Started consuming from topic: %s", c.topicName)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var evt event.InventoryEvent
			if err := json.Unmarshal(msg.Value, &evt); err != nil {
				log.Printf("Error unmarshaling event: %v", err)
				continue
			}

			if err := c.processor.ProcessInventoryEvent(ctx, evt); err != nil {
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
