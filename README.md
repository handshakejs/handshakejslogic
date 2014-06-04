# handshakejslogic

<img src="http://i.imgur.com/IlE7oAi.png" alt="handshakejslogic" align="right" />

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

## Installation

```
go get github.com/handshakejs/handshakejslogic
```

## Running Tests

```
go test -v
```
