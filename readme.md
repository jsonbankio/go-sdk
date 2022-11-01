# Jsonbank.io GoLang Sdk

The official repository for jsonbank.io GoLang SDK.

## Installation

```bash
go get github.com/jsonbankio/go-sdk
```

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/jsonbankio/go-sdk"
)

func main() {
	// Initialize the client
	jsb := jsonbank.InitWithoutKeys()

	// Get Public content
	content, err := jsb.GetContent("jsonbank/sdk-test/index.json")
	if err != nil {
		panic(err)
	}

	// convert to json string
	jsonData, _ := json.MarshalIndent(content, "", "  ")

	fmt.Println(string(jsonData))
}
```

### Authenticated Requests

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/jsonbankio/go-sdk"
	"github.com/jsonbankio/go-sdk/types"
)

func main() {
	// Initialize the client
	jsb := jsonbank.Init(jsonbank.Config{
		Keys: jsonbank.Keys{
			Public:  "your public key",
			Private: "your private key",
		},
	})

	// Authenticate to check if the keys are valid
	authenticate, err := jsb.Authenticate()
	if err != nil {
		fmt.Println(err)
	}

	// Print authenticated user
	fmt.Println("Authenticated as: ", authenticate.Username)

	// Get Own content
	content, err := jsb.GetOwnContent("sdk-test/index.json")
	if err != nil {
		panic(err)
	}

	// convert to json string
	jsonData, _ := json.MarshalIndent(content, "", "  ")

	fmt.Println(string(jsonData))
}
```

### Testing

Create an .env file in the root of the project and add the following variables

```dotenv
JSB_HOST="https://api.jsonbank.io"
JSB_PUBLIC_KEY="your public key"
JSB_PRIVATE_KEY="your private key"
```

Then run the tests

```bash
go test -v
```