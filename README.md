# handshakejslogic

<img src="https://raw.githubusercontent.com/handshakejs/handshakejslogic/master/handshakejslogic.gif" alt="handshakejslogic" align="right" width="190" />

[![BuildStatus](https://travis-ci.org/handshakejs/handshakejslogic.png?branch=master)](https://travis-ci.org/handshakejs/handshakejslogic)

Logic for saving handshakejs data to the redis database.

This library is part of the larger [Handshake.js ecosystem](https://github.com/handshakejs).

## Usage

```go
package main

import (
  "fmt"
  handshakejslogic "github.com/handshakejs/handshakejslogic"
)

func main() {
  handshakejslogic.Setup("redis://127.0.0.1:6379")

  app := map[string]interface{}{"email": "email@myapp.com", "app_name": "myapp"}
  result, logic_error := handshakejslogic.AppsCreate(app)
  if logic_error != nil {
    fmt.Println(logic_error)
  }
  fmt.Println(result)
```

### Setup

Connects to Redis.

```go
handshakejslogic.Setup("redis://127.0.0.1.6379")
```

### AppsCreate

```go
app := map[string]interface{}{"email": "email@myapp.com", "app_name": "myapp"}
result, logic_error := handshakejslogic.AppsCreate(app)
```

### IdentitiesCreate

```go
identity := map[string]interface{}{"email": "user@email.com", "app_name": "myapp"}
result, logic_error := handshakejslogic.IdentitiesCreate(identity)
```

## Installation

```
go get github.com/handshakejs/handshakejslogic
```

## Running Tests

```
go test -v
```
