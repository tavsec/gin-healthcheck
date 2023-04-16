# Gin Healthcheck
[![Go Reference](https://pkg.go.dev/badge/github.com/tavsec/gin-healthcheck.svg)](https://pkg.go.dev/github.com/tavsec/gin-healthcheck)
![tests](https://github.com/tavsec/gin-healthcheck/actions/workflows/test.yaml/badge.svg)

This module will create a simple endpoint for Gin framework, 
which can be used to determined healthiness of Gin application.

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
)

func main() {
    r := gin.Default()

    healthcheck.New(r, healthcheck.DefaultConfig(), []checks.Check{})
	
    r.Run()
}
```

This will add the healthcheck endpoint to default path, which is `/healthz`. The path can be customized
using `healthcheck.Config` structure. In the example above, no specific checks will be included, only API availability.

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
In case you want to ensure that your application can reach seperate service, 
you can utilise `PingCheck`.

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

### Environmental variables check
You can check if environmental variable is set using `EnvCheck`:
```go
package main

import (
	"github.com/gin-gonic/gin"
    healthcheck "github.com/tavsec/gin-healthcheck"
    "github.com/tavsec/gin-healthcheck/checks"
    "github.com/tavsec/gin-healthcheck/config"
)

func main(){
    r := gin.Default()

    dbHostCheck := checks.NewEnvCheck("DB_HOST")
	
	// You can also validate env format using regex
    dbUserCheck := checks.NewEnvCheck("DB_HOST")
	dbUserCheck.SetRegexValidator("^USER_")
	
    healthcheck.New(r, config.DefaultConfig(), []checks.Check{dbHostCheck, dbUserCheck})

    r.Run()
}

```

### Custom checks
Besides built-in health checks, you can extend the functionality and create your own check, utilising the `Check` interface: 
```go
package checks

type Check interface {
    Pass() bool
    Name() string
}
```
