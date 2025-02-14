package rabbitmq

import (
	"fmt"

	"github.com/org/2112-space-lab/org/app-service/internal/config"
	log "github.com/org/2112-space-lab/org/app-service/pkg/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient wraps the RabbitMQ connection and channel.
type RabbitMQClient struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	env         *config.SEnv
	outputQueue string
	inputQueue  string
}

// NewRabbitMQClient initializes a new RabbitMQ client.
func NewRabbitMQClient(env *config.SEnv) (*RabbitMQClient, error) {
	// Connect to RabbitMQ
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
		outputQueue: env.EnvVars.RabbitMQ.OutputQueue, // Get output queue name from config
		inputQueue:  env.EnvVars.RabbitMQ.InputQueue,  // Get input queue name from config
	}

	// Setup queues
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

// SetupQueues ensures both input and output queues exist and are bound correctly.
func (r *RabbitMQClient) SetupQueues() error {
	exchangeName := "app.exchange" // Default exchange name, replace if needed

	// Declare and bind input queue
	if _, err := r.DeclareQueue(r.inputQueue); err != nil {
		return fmt.Errorf("failed to declare input queue: %w", err)
	}
	if err := r.BindQueue(r.inputQueue, r.inputQueue, exchangeName); err != nil {
		return fmt.Errorf("failed to bind input queue: %w", err)
	}

	// Declare and bind output queue
	if _, err := r.DeclareQueue(r.outputQueue); err != nil {
		return fmt.Errorf("failed to declare output queue: %w", err)
	}
	if err := r.BindQueue(r.outputQueue, r.outputQueue, exchangeName); err != nil {
		return fmt.Errorf("failed to bind output queue: %w", err)
	}

	log.Debugf("✅ Queues setup complete: Input=%s, Output=%s", r.inputQueue, r.outputQueue)
	return nil
}

// DeclareExchange declares an exchange in RabbitMQ.
func (r *RabbitMQClient) DeclareExchange(exchangeName, exchangeType string) error {
	return r.channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
}

// DeclareQueue declares a queue in RabbitMQ.
func (r *RabbitMQClient) DeclareQueue(queueName string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
}

// BindQueue binds a queue to an exchange with a routing key.
func (r *RabbitMQClient) BindQueue(queueName, routingKey, exchangeName string) error {
	return r.channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}

// PublishMessage sends a message to the configured output queue.
func (r *RabbitMQClient) PublishMessage(body []byte) error {
	err := r.channel.Publish(
		"",
		r.outputQueue, // Use the output queue from config
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to output queue: %w", err)
	}

	log.Debugf("[x] Message published to output queue: %s", r.outputQueue)
	return nil
}

// ConsumeMessages listens for messages from the configured input queue.
func (r *RabbitMQClient) ConsumeMessages() (<-chan amqp.Delivery, error) {
	msgs, err := r.channel.Consume(
		r.inputQueue, // Use the input queue from config
		"",
		true,  // Auto-Acknowledge
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer on input queue: %w", err)
	}
	log.Debugf("📥 Listening for messages on input queue: %s", r.inputQueue)
	return msgs, nil
}
