package rabbitmq

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn *amqp091.Connection
	ch *amqp091.Channel
}

// creates new connection and channel
func New(url string) (*Client, error) {
	conn, err := amqp091.Dial(url) // new connection
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open rabbitmq channel: %w", err)
	}

	return &Client{
		conn: conn,
		ch: ch,
	}, nil
}

// creates new queue
func (c *Client) DeclareQueue(name string) error {
	_, err := c.ch.QueueDeclare(
		name, 
		true,
		false,
		false,
		false,
		nil, 
	)

	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", name, err)
	}

	return nil
}

// message publishing to an exact queue
func (c *Client) Publish(ctx context.Context, queueName string, body []byte) error {
	err := c.ch.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			DeliveryMode: amqp091.Persistent,
			Body: body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message to %s: %w", queueName, err)
	}

	return nil
}

// getting signed up to the queue (Ack confirms succsessful msg processing)
func (c *Client) Consume(queueName string) (<-chan amqp091.Delivery, error) {
	deliveries, err := c.ch.Consume(
		queueName,
		"",
		false, //AutoAck (rabbit wont delete msq from queue till worker havent sent msq.Ack(false))
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to consume queue %s: %w", queueName, err)
	}

	return deliveries, nil
}

func (c *Client) Close() error {
	if c.ch != nil {
		_ = c.ch.Close()
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}