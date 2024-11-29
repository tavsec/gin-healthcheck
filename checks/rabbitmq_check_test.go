package checks

import (
	"context"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

func TestRabbit_Check(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Setup the RabbitMQ container
	rabbitmqContainer, err := rabbitmq.Run(ctx, "rabbitmq:3.12.11-management-alpine")
	require.NoError(t, err)

	// Ensure the container is cleaned up
	t.Cleanup(func() {
		_ = rabbitmqContainer.Terminate(ctx)
	})

	url, err := rabbitmqContainer.AmqpURL(ctx)
	require.NoError(t, err)

	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Connection is nil",
			want: false,
		},
		{
			name: "Connection is not closed",
			want: true,
		},
		{
			name: "Connection is closed",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var conn *amqp.Connection
			switch tt.name {
			case "Connection is closed":
				conn, err = amqp.Dial(url)
				require.NoError(t, err)

				time.Sleep(time.Second)

				err = conn.Close()
				require.NoError(t, err)

			case "Connection is not closed":
				conn, err = amqp.Dial(url)
				require.NoError(t, err)
			}

			r := NewRabbitCheck(conn)
			assert.Equal(t, tt.want, r.Pass())
		})
	}
}

func TestRabbitMQ_Name(t *testing.T) {
	r := &RabbitMQCheck{}
	if got := r.Name(); got != "rabbitmq" {
		t.Errorf("RabbitMQCheck.Name() = %v, want %v", got, "rabbitmq")
	}
}
