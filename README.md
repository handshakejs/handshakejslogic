<img src="https://raw.githubusercontent.com/handshakejs/handshakejslogic/master/handshakejslogic.gif" alt="handshakejslogic" align="right" />

# handshakejslogic

Logic for saving handshakejs data to the redis database.

[![BuildStatus](https://travis-ci.org/handshakejs/handshakejslogic.png?branch=master)](https://travis-ci.org/handshakejs/handshakejslogic)

## Installation

```
go get github.com/handshakejs/handshakejslogic
```

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
  result, logic_error := handshakejslogic.AppsCreate(app)
  if logic_error != nil {
    fmt.Println(logic_error)
  }
  fmt.Println(result)

  identity := map[string]interface{}{"email": "identity0@mailinator.com", "app_name": APP_NAME}
  result2, logic_error := handshakejslogic.IdentitiesCreate(app)
  if logic_error != nil {
    fmt.Println(logic_error)
  }
  fmt.Println(result2)
}
```

## Running Tests

```
go test -v
```
