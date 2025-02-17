package rabbitmq

import (
	"fmt"

	"github.com/org/2112-space-lab/org/app-service/internal/config"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExhangeDefaultName     = "headers.exchange"
	ExchangeTypeHeaders    = "headers"
	RabbitMQFormat         = "application/json"
	DefaultAutoAcknowledge = true
	DefaultAutoDelete      = false
	DefaultNotExclusive    = false
	DefaultNoLocal         = false
	DefaultNoWait          = false
	DefaultDurable         = true
	DefaultImmediate       = false
	DefaultMandatory       = false
	DefaultKey             = ""
)

// RabbitMQClient wraps the RabbitMQ connection and channel.
type RabbitMQClient struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	env         *config.SEnv
	outputQueue string
	inputQueue  string
	exchange    string
	defaultArgs amqp.Table
}

// NewRabbitMQClient initializes a new RabbitMQ client.
func NewRabbitMQClient(env *config.SEnv) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(env.EnvVars.RabbitMQ.GetAddr())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	client := &RabbitMQClient{
		conn:        conn,
		channel:     ch,
		env:         env,
		outputQueue: env.EnvVars.RabbitMQ.OutputQueue,
		inputQueue:  env.EnvVars.RabbitMQ.InputQueue,
		exchange:    ExhangeDefaultName,
		defaultArgs: amqp.Table{},
	}

	if err := client.SetupQueues(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to setup queues: %w", err)
	}

	return client, nil
}

// Close closes the RabbitMQ connection and channel.
func (r *RabbitMQClient) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// SetupQueues ensures queues and the exchange exist and are bound correctly.
func (r *RabbitMQClient) SetupQueues() error {
	if err := r.DeclareExchange(r.exchange, ExchangeTypeHeaders); err != nil {
		return fmt.Errorf("failed to declare headers exchange: %w", err)
	}

	if _, err := r.DeclareQueue(r.inputQueue); err != nil {
		return fmt.Errorf("failed to declare input queue: %w", err)
	}

	if _, err := r.DeclareQueue(r.outputQueue); err != nil {
		return fmt.Errorf("failed to declare output queue: %w", err)
	}

	log.Debugf("✅ Queues and exchange setup complete: Exchange=%s, Input=%s, Output=%s", r.exchange, r.inputQueue, r.outputQueue)
	return nil
}

// DeclareExchange declares a headers exchange in RabbitMQ.
func (r *RabbitMQClient) DeclareExchange(exchangeName, exchangeType string) error {
	return r.channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		DefaultDurable,
		DefaultAutoDelete,
		DefaultNoLocal,
		DefaultNoWait,
		r.defaultArgs,
	)
}

// DeclareQueue declares a queue in RabbitMQ.
func (r *RabbitMQClient) DeclareQueue(queueName string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		queueName,
		DefaultDurable,
		DefaultAutoDelete,
		DefaultNoLocal,
		DefaultNoWait,
		r.defaultArgs,
	)
}

// PublishMessage sends a message with dynamic headers.
func (r *RabbitMQClient) PublishMessage(body []byte, headers *Header) error {
	err := r.channel.Publish(
		r.exchange,
		DefaultKey,
		DefaultMandatory,
		DefaultImmediate,
		amqp.Publishing{
			ContentType: RabbitMQFormat,
			Body:        body,
			Headers:     headers.Fields,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to exchange: %w", err)
	}

	log.Debugf("[x] Message published with headers: %+v", headers.Fields)
	return nil
}

// ConsumeMessages listens for messages that match specific headers.
func (r *RabbitMQClient) ConsumeMessages(filterHeaders *Header) (<-chan amqp.Delivery, error) {
	_, err := r.DeclareQueue(r.inputQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	err = r.channel.QueueBind(
		r.inputQueue,
		DefaultKey,
		r.exchange,
		DefaultNoWait,
		amqp.Table(filterHeaders.Fields),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue with headers: %w", err)
	}

	msgs, err := r.channel.Consume(
		r.inputQueue,
		DefaultKey,
		DefaultAutoAcknowledge,
		DefaultNotExclusive,
		DefaultNoLocal,
		DefaultNoWait,
		r.defaultArgs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Debugf("📥 Listening for messages with filters: %+v", filterHeaders.Fields)
	return msgs, nil
}
