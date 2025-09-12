package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"rest-api-in-gin/internal/payment/domain/service"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// UserCreatedEvent is a struct that the consumer expects to receive
type UserCreatedEvent struct {
	ID string `json:"id"`
}

// UserCreatedConsumer is a struct to listen to the Outbox database table and a Kafka consumer
type UserCreatedConsumer struct {
	reader              *kafka.Reader
	walletDomainService *service.WalletDomainService
}

// NewUserCreatedConsumer creates a new instance of UserCreatedConsumer
func NewUserCreatedConsumer(kafkaReader *kafka.Reader, walletDomainService *service.WalletDomainService) *UserCreatedConsumer {
	return &UserCreatedConsumer{
		reader:              kafkaReader,
		walletDomainService: walletDomainService,
	}
}

// Run starts the continuous listening process
func (c *UserCreatedConsumer) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("UserCreatedConsumer stopping")
			return
		default:
			c.processMessage(ctx)
			time.Sleep(1 * time.Second) // optional small delay
		}
	}
}

// processMessages consumes new events from the Kafka broker and process them
func (c *UserCreatedConsumer) processMessage(ctx context.Context) {
	m, err := c.reader.ReadMessage(ctx)
	if err != nil {
		if ctx.Err() != nil {
			log.Println("context cancelled, stopping consumer")
			return
		}
		log.Println("error reading from kafka:", err)
		return
	}

	var evt UserCreatedEvent
	if err := json.Unmarshal(m.Value, &evt); err != nil {
		log.Println("invalid event payload:", err)
		return
	}

	userID, err := uuid.Parse(evt.ID)
	if err != nil {
		log.Println("invalid UUID in event:", evt.ID)
		return
	}

	if err := c.walletDomainService.CreateWalletForNewUser(userID); err != nil {
		log.Println("failed to create wallet:", err)
	} else {
		log.Printf("wallet created for user %s\n", evt.ID)
	}
}
