package kafka

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/knoxgao67/temporal-dispatcher-demo/config"
)

const (
	consumerGroup = "temporal-demo-group"
	consumerTopic = "temporal-demo-topic"
)

var consumer *Consumer

func initConsumer() {
	var err error
	consumer, err = NewConsumer(config.GetConfig().KafkaBrokers, consumerTopic)
	if err != nil {
		log.Fatalln("Unable to create Kafka consumer.", err)
	}
}

func GetConsumer() *Consumer {
	return consumer
}

type Consumer struct {
	underlying sarama.ConsumerGroup
	msgChan    chan *sarama.ConsumerMessage
}

func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()} // default rr
	underlying, err := sarama.NewConsumerGroup(brokers, consumerGroup, cfg)
	if err != nil {
		return nil, err
	}

	c := &Consumer{
		underlying: underlying,
		msgChan:    make(chan *sarama.ConsumerMessage),
	}

	go c.StartConsumer(context.Background(), topic)

	return c, nil
}

func (c *Consumer) StartConsumer(ctx context.Context, topic string) {
	handler := &consumerGroupHandler{
		msgChan: c.msgChan,
		ready:   make(chan bool),
	}
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := c.underlying.Consume(ctx, []string{topic}, handler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			}
			log.Printf("Error from consumer: %v\n", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}
		handler.ready = make(chan bool)
	}
}

func (c *Consumer) ReadMessage(duration time.Duration) string {
	select {
	case msg, ok := <-c.msgChan:
		if !ok {
			return ""
		}
		return string(msg.Value)
	case <-time.NewTimer(duration).C:
		return ""
	}
}

func (c *Consumer) Close() error {
	close(c.msgChan)
	if err := c.underlying.Close(); err != nil {
		log.Printf("Failed to close underlying: %v", err)
		return err
	}

	return nil
}

type consumerGroupHandler struct {
	ready   chan bool
	msgChan chan *sarama.ConsumerMessage
}

func (c *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}
func (c *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("chan is close,return")
		}
	}()
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("msgChan channel was closed")
				return nil
			}
			// notice may be panic, use recover to catch it
			c.msgChan <- message
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
