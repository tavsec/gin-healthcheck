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
)

func main() {
    r := gin.Default()

    healthcheck.New(r, healthcheck.DefaultConfig())
	
    r.Run()
}
```

This will add the healthcheck endpoint to default path, which is `/healthz`. The path can be customized
using `healthcheck.Config` structure.
