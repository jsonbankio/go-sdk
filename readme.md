# Jsonbank.io GoLang Sdk

##### Still in Development Phase

## Usage

```go
package main

import (
	"fmt"
	"github.com/jsonbankio/go-sdk"
	"github.com/jsonbankio/go-sdk/types"
)

func main() {
	jsb := jsonbank.Init(jsonbank.Config{
		Keys: jsonbank.Keys{
			Public:  "your public key",
			Private: "your private key",
		},
	})

	authenticate, err := jsb.Authenticate()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Authenticated as: ", authenticate.Username)

	meta, err := jsb.GetOwnDocumentMeta("js-sdk-test/index.json")
	if dataError != nil {
		fmt.Println("Error:", dataError)
		return
	}

	fmt.Println("Data:", meta.Path)

	const testFileContent = `{
    "name": "Js SDK Test File",
    "author": "jsonbank"
	}`

	document, err := jsb.CreateDocumentIfNotExists(types.CreateDocumentBody{
		Name:    "index.json",
		Content: testFileContent,
		Project: "js-sdk-test",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Document:", document)
}
```