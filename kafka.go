package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer wraps a Kafka writer for sending analytics events
type KafkaProducer struct {
	writer  *kafka.Writer
	enabled bool
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(config *Config) *KafkaProducer {
	if !config.KafkaEnabled {
		log.Println("Kafka is disabled, events will only be written to file")
		return &KafkaProducer{enabled: false}
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.KafkaBrokers...),
		Topic:        config.KafkaTopic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    1,
		BatchTimeout: 10 * time.Millisecond,
		Async:        true, // Fire and forget for analytics
	}

	log.Printf("Kafka producer initialized for topic %s (brokers: %v)", config.KafkaTopic, config.KafkaBrokers)
	return &KafkaProducer{
		writer:  writer,
		enabled: true,
	}
}

// SendEvent sends an event to Kafka
func (kp *KafkaProducer) SendEvent(event map[string]interface{}) error {
	if !kp.enabled {
		return nil
	}

	// Add timestamp if not present
	if _, ok := event["timestamp"]; !ok {
		event["timestamp"] = time.Now().Format(time.RFC3339)
	}

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return err
	}

	msg := kafka.Message{
		Key:   []byte(event["type"].(string)),
		Value: data,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = kp.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Failed to send event to Kafka: %v", err)
		return err
	}

	return nil
}

// Close closes the Kafka writer
func (kp *KafkaProducer) Close() error {
	if kp.enabled && kp.writer != nil {
		return kp.writer.Close()
	}
	return nil
}

