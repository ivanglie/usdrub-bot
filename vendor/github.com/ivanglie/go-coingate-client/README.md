# Golang client for CoinGate exchange rate API

[![Go Reference](https://pkg.go.dev/badge/github.com/ivanglie/go-coingate-client.svg)](https://pkg.go.dev/github.com/ivanglie/go-coingate-client)
[![Test](https://github.com/ivanglie/go-coingate-client/actions/workflows/test.yml/badge.svg)](https://github.com/ivanglie/go-coingate-client/actions/workflows/test.yml)

Golang client for the [CoinGate Exchange Rate API](https://developer.coingate.com/docs/get-rate).

## Example

First, ensure the library is installed and up to date by running go get -u github.com/ivanglie/go-coingate-client.

This is a very simple app that just displays US Dollar to Chinese Yuan Renminbi conversion.

```golang
package main

import (
	"fmt"
	"time"

	coingate "github.com/ivanglie/go-coingate-client"
)

func main() {
	client := coingate.NewClient()
	rate, err := client.GetRate("USD", "CNY")
	if err != nil {
		panic(err)
	}
	fmt.Println(rate)
}
```
See [main.go](./_example/main.go).

## References

For more information check out the following links:

* [CoinGate Exchange Rate API](https://developer.coingate.com/docs/get-rate)