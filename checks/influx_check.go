package checks

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxV2Check struct {
	client  influxdb2.Client
	Timeout int
}

func NewInfluxV2Check(timeout int, client influxdb2.Client) *InfluxV2Check {
	return &InfluxV2Check{
		client:  client,
		Timeout: timeout,
	}
}

func (i *InfluxV2Check) Pass() bool {
	if i.client == nil {
		return false
	}

	timeout := time.Second * time.Duration(i.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ping, err := i.client.Ping(ctx)
	if err != nil {
		return false
	}

	return ping
}

func (i *InfluxV2Check) Name() string {
	return "influxdb"
}
