# Golang client for exchange rate of cash currency in Russia

[![Go Reference](https://pkg.go.dev/badge/github.com/ivanglie/go-br-client.svg)](https://pkg.go.dev/github.com/ivanglie/go-br-client)
[![Test](https://github.com/ivanglie/go-br-client/actions/workflows/test.yml/badge.svg)](https://github.com/ivanglie/go-br-client/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/ivanglie/go-br-client/branch/master/graph/badge.svg?token=8lRyze5RSQ)](https://codecov.io/gh/ivanglie/go-br-client)

Golang client that provides latest exchange rate of USDRUB cash in largest cities of Russia.

## Example

First, ensure the library is installed and up to date by running

```
go get -u github.com/ivanglie/go-br-client
```

This is a very simple app that just displays USDRUB exhange rate in Novosibirsk.

```golang
package main

import (
	"fmt"

	br "github.com/ivanglie/go-br-client"
)

func main() {
	client := br.NewClient()
	rates, err := client.Rates(br.Novosibirsk)
	if err != nil {
		panic(err)
	}
	fmt.Println(rates)
}
```

Console output:

```json
{
    "currency": "USD",
    "city": "novosibirsk",
    "branches": [
        {
            "bank": "ОП № 029/0000 Филиала \"Газпромбанк\" АО",
            "subway": "Заельцовская, Гагаринская, Сибирская",
            "currency": "USD",
            "buy": 87.1,
            "sell": 91.3,
            "updated": "2023-07-03T13:00:00+03:00"
        },
        {
            "bank": "ДО № 029/1007 Филиала \"Газпромбанк\" АО",
            "subway": "Заельцовская, Берёзовая роща, Гагаринская",
            "currency": "USD",
            "buy": 87.1,
            "sell": 91.3,
            "updated": "2023-07-03T13:00:00+03:00"
        },
        {
            "bank": "ДО № 029/1003 Филиала \"Газпромбанк\" АО",
            "subway": "Берёзовая роща, Маршала Покрышкина, Золотая Нива",
            "currency": "USD",
            "buy": 87.1,
            "sell": 91.3,
            "updated": "2023-07-03T13:00:00+03:00"
        }
    ]
}
```
See [main.go](./_example/main.go).

## References

For more information check out the following links:

* Cash currency exchange rates by [Banki.ru](https://www.banki.ru/products/currency/map/moskva/) (RU)
