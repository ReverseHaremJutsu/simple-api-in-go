package producer

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// OutboxEvent is a struct to hold the rows queried before publishing to Kafka broker
type OutboxEvent struct {
	ID            uuid.UUID
	AggregateName string
	AggregateID   uuid.UUID
	EventName     string
	Payload       json.RawMessage
	TopicName     string
	CreatedAt     time.Time
}

// OutboxDelivery is a struct to listen to the Outbox database table and a Kafka producer
type OutboxDelivery struct {
	db          *sql.DB
	kafkaWriter *kafka.Writer
}

// NewOutboxDelivery creates a new instance of OutboxDelivery
func NewOutboxDelivery(db *sql.DB, kafkaWriter *kafka.Writer) *OutboxDelivery {
	return &OutboxDelivery{db: db, kafkaWriter: kafkaWriter}
}

// Run starts the continuous OutboxDelivery process
func (d *OutboxDelivery) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			d.processEvents(ctx)
			time.Sleep(30 * time.Second) // poll every 30s
		}
	}
}

// processEvents polls for new events from the database and publishes to Kafka
func (d *OutboxDelivery) processEvents(ctx context.Context) {
	rows, err := d.db.QueryContext(ctx, `
		SELECT id, aggregate_name, aggregate_id, event_name, payload, topic_name, created_at
		FROM outbox
		ORDER BY created_at ASC
		LIMIT 10
	`)
	if err != nil {
		log.Println("failed to fetch outbox:", err)
		return
	}
	defer rows.Close()

	var events []OutboxEvent
	for rows.Next() {
		var e OutboxEvent
		if err := rows.Scan(&e.ID, &e.AggregateName, &e.AggregateID, &e.EventName, &e.Payload, &e.TopicName, &e.CreatedAt); err != nil {
			log.Println("failed scan:", err)
			continue
		}
		events = append(events, e)
	}

	for _, e := range events {
		err := d.kafkaWriter.WriteMessages(ctx,
			kafka.Message{
				Key:   []byte(e.AggregateID.String()),
				Value: e.Payload,
				Topic: e.TopicName,
			})
		if err != nil {
			log.Println("failed publish to kafka:", err)
			continue
		}

		_, err = d.db.ExecContext(ctx, `DELETE FROM outbox WHERE id = $1`, e.ID)
		if err != nil {
			log.Println("failed delete outbox:", err)
		}
	}
}
