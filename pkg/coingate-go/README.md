# Golang client for CoinGate exchange rate API

Golang client for the [CoinGate Exchange](https://developer.coingate.com/docs/get-rate).

## Example

First, ensure the library is installed and up to date by running

```
go get -u github.com/ivanglie/usdrub-bot/pkg/coingate-go
```

This is a very simple app that just displays US Dollar to Chinese Yuan Renminbi conversion.

```golang
package main

import (
	"fmt"

	"github.com/ivanglie/usdrub-bot/pkg/coingate-go"
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
See [main.go](../../examples/coingate/main.go).

## References

For more information check out the following links:

* [CoinGate Exchange Rate API](https://developer.coingate.com/docs/get-rate)