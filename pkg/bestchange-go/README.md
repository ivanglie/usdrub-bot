# Golang client for exchange rate of USDT (TRC20) to RUB cash in Russia

Golang client that provides latest exchange rate of USDRUB cash in largest cities of Russia.

## Example

First, ensure the library is installed and up to date by running

```
go get -u github.com/ivanglie/usdrub-bot/pkg/bestchange-go
```

This is a very simple app that just displays USDRUB exhange rate in Novosibirsk.

```golang
package main

import (
	"fmt"

	"github.com/ivanglie/usdrub-bot/pkg/bestchange-go"
)

func main() {
	client := bestchange.NewClient()
	rate, err := client.Rate()
	if err != nil {
		panic(err)
	}
	fmt.Println(rate)
}
```

Console output:

```
95.920519
```
See [main.go](../../examples/crypto/main.go).

## References

For more information check out the following links:

* Cash currency exchange rates by [Banki.ru](https://www.banki.ru/products/currency/map/moskva/) (RU)
