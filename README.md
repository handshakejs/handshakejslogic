# handshakejslogic

Logic for saving handshakejs data to the redis database.

[![BuildStatus](https://travis-ci.org/handshakejs/handshakejslogic.png?branch=master)](https://travis-ci.org/handshakejs/handshakejslogic)

## Usage

```go
package main

import (
  "fmt"
  handshakejslogic "github.com/handshakejs/handshakejslogic"
)

func main() {
  handshakejslogic.Setup("redis://127.0.0.1:6379")

  app := map[string]interface{}{"email": EMAIL, "app_name": APP_NAME}
  result, err := handshakejslogic.AppsCreate(app)
  if err != nil {
    fmt.Println(err)
  }

  fmt.Println(result)
}
```

## Installation

```
go get github.com/handshakejs/handshakejslogic
```

## Running Tests

```
go test -v
```
