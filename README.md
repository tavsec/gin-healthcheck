# Gin Healthcheck
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

### Adding additional checks
Currently, gin-healthcheck comes with SQL check, which will send `ping` request to SQL.
```go
package main

import (
    "github.com/gin-gonic/gin"
    healthcheck "github.com/tavsec/gin-healthcheck"
    "github.com/tavsec/gin-healthcheck/checks"
)

func main() {
    r := gin.Default()

    // Initialize Database
    // db := ...
    // ...
    healthcheck.New(r, healthcheck.DefaultConfig(), []checks.Check{checks.SqlCheck{Sql: db}})
	
    r.Run()
}
```

Besides that, you can extend the functionality and create your own check, utilising the `Check` interface: 
```go
package checks

type Check interface {
    Pass() bool
    Name() string
}
```
