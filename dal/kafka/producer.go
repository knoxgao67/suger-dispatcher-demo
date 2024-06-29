package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/knoxgao67/temporal-dispatcher-demo/config"
)

var producer *Producer

func initProducer() {
	var err error
	producer, err = NewProducer(config.GetConfig().KafkaBrokers, consumerTopic)
	if err != nil {
		log.Fatalln("Unable to create Kafka producer.", err)
	}
}

func GetProducer() *Producer {
	return producer
}

type Producer struct {
	topic      string
	underlying sarama.SyncProducer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner // default rr

	underlying, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return &Producer{
		topic:      topic,
		underlying: underlying,
	}, nil
}

func (p *Producer) SendMessage(ctx context.Context, msg string) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(msg),
	}

	_, _, err := p.underlying.SendMessage(message)
	if err != nil {
		log.Printf("Failed to send message to topic %s: %v", p.topic, err)
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.underlying.Close(); err != nil {
		log.Printf("Failed to close underlying: %v", err)
		return err
	}

	return nil
}
