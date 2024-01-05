# Gin Healthcheck

[![Go Reference](https://pkg.go.dev/badge/github.com/tavsec/gin-healthcheck.svg)](https://pkg.go.dev/github.com/tavsec/gin-healthcheck)
![tests](https://github.com/tavsec/gin-healthcheck/actions/workflows/test.yaml/badge.svg)

This module will create a simple endpoint for Gin framework,
which can be used to determine the healthiness of Gin application.

## Installation

Install package:

```shell
go get github.com/tavsec/gin-healthcheck
```

## Usage

```go
package main

import (
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

func main() {
	r := gin.Default()

	healthcheck.New(r, config.DefaultConfig(), []checks.Check{})

	r.Run()
}
```

This will add the healthcheck endpoint to the default path, which is `/healthz`. The path can be customized
using `config.Config` structure. In the example above, no specific checks will be included, only API availability.

## Health checks

### SQL

Currently, gin-healthcheck comes with SQL check, which will send `ping` request to SQL.

```go
package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

func main() {
	r := gin.Default()

	db, _ := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/hello")
	sqlCheck := checks.SqlCheck{Sql: db}
	healthcheck.New(r, config.DefaultConfig(), []checks.Check{sqlCheck})

	r.Run()
}
```

### Ping

In case you want to ensure that your application can reach a separate service, you can utilize `PingCheck`.

```go
package main

import (
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

func main() {
	r := gin.Default()

	pingCheck := checks.NewPingCheck("https://www.google.com", "GET", 1000, nil, nil)
	healthcheck.New(r, config.DefaultConfig(), []checks.Check{pingCheck})

	r.Run()
```

### Redis check

You can perform Redis ping check using `RedisCheck` checker:

```go
package main

import (
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
	"github.com/redis/go-redis/v9"
)

func main() {
	r := gin.Default()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redisCheck := checks.NewRedisCheck(rdb)
	healthcheck.New(r, config.DefaultConfig(), []checks.Check{redisCheck})

	r.Run()
```

### MongoDB check

You can perform MongoDB ping check using `Mongo` checker:

```go
package main

import (
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := gin.Default()

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, _ := mongo.Connect(context.TODO(), clientOptions)

	mongoCheck := checks.NewMongoCheck(10, client)
	healthcheck.New(r, config.DefaultConfig(), []checks.Check{mongoCheck})

	r.Run()
}
```

### InfluxDB V2 check

You can perform InfluxDB V2 ping check using `InfluxV2` checker:

```go
package main

import (
	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

func main() {
	r := gin.Default()

	// Create an influxdbv2 client
	client := influxdb2.NewClient("localhost:8086", "token")

	influxCheck := checks.NewInfluxV2Check(10, client)
	healthcheck.New(r, config.DefaultConfig(), []checks.Check{influxCheck})

	r.Run()
}

```

### Environmental variables check

You can check if an environmental variable is set using `EnvCheck`:

```go
package main

import (
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

func main() {
	r := gin.Default()

	dbHostCheck := checks.NewEnvCheck("DB_HOST")

	// You can also validate env format using regex
	dbUserCheck := checks.NewEnvCheck("DB_HOST")
	dbUserCheck.SetRegexValidator("^USER_")

	healthcheck.New(r, config.DefaultConfig(), []checks.Check{dbHostCheck, dbUserCheck})

	r.Run()
}
```

### context.Context check

You can check if a context has not been canceled, by using a `ContextCheck`:

```go
package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

func main() {
	r := gin.Default()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	signalsCheck := checks.NewContextCheck(ctx, "signals")
	healthcheck.New(r, config.DefaultConfig(), []checks.Check{signalsCheck})

	r.Run()
}
```

### Custom checks

Besides built-in health checks, you can extend the functionality and create your own check, utilizing the `Check`
interface:

```go
package checks

type Check interface {
	Pass() bool
	Name() string
}
```

## Notification of health check failure

It is possible to get notified when the health check failed a certain threshold of call. This would match for example
the failureThreshold of Kubernetes and allow us to take action in that case.

```go
package main

import (
	"github.com/gin-gonic/gin"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
)

func main() {
	r := gin.Default()

	conf := healthcheck.DefaultConfig()

	conf.FailureNotification.Chan = make(chan error, 1)
	defer close(conf.FailureNotification.Chan)
	conf.FailureNotification.Threshold = 3

	go func() {
		<-conf.FailureNotification.Chan
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	signalsCheck := checks.NewContextCheck(ctx, "signals")
	healthcheck.New(r, conf, []checks.Check{signalsCheck})

	r.Run()
}
```

Note that the following example is not doing a graceful shutdown. If Kubernetes is set up with a failureThreshold of 3,
it will mark the pod as failing after that third call, but there is no guarantee that you have processed and answered
all HTTP requests before the call to os.Exit(1). It is necessary to use something
like https://github.com/gin-contrib/graceful at that point to have a graceful shutdown.
