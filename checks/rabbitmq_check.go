package checks

import amqp "github.com/rabbitmq/amqp091-go"

type RabbitMQCheck struct {
	conn *amqp.Connection
}

func NewRabbitCheck(conn *amqp.Connection) *RabbitMQCheck {
	return &RabbitMQCheck{
		conn: conn,
	}
}

func (r *RabbitMQCheck) Pass() bool {
	if r.conn == nil {
		return false
	}

	return !r.conn.IsClosed()
}

func (r *RabbitMQCheck) Name() string {
	return "rabbitmq"
}
