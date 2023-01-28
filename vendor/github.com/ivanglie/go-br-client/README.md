# Golang client for exchange rate of cash currency in Russia

[![Go Reference](https://pkg.go.dev/badge/github.com/ivanglie/go-br-client.svg)](https://pkg.go.dev/github.com/ivanglie/go-br-client)
[![Test](https://github.com/ivanglie/go-br-client/actions/workflows/test.yml/badge.svg)](https://github.com/ivanglie/go-br-client/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/ivanglie/go-br-client/branch/master/graph/badge.svg?token=8lRyze5RSQ)](https://codecov.io/gh/ivanglie/go-br-client)

Golang client that provides latest exchange rate of cash currency in largest cities of Russia.

## Example

First, ensure the library is installed and up to date by running

```
go get -u github.com/ivanglie/go-br-client
```

This is a very simple app that just displays exhange rate of Chinese Yuan Renminbi in Novosibirsk.

```golang
package main

import (
	"fmt"

	br "github.com/ivanglie/go-br-client"
)

func main() {
	client := br.NewClient()
	rates, err := client.Rates(br.CNY, br.Novosibirsk)
	if err != nil {
		panic(err)
	}
	fmt.Println(rates)
}
```

Console output:

```json
{
    "currency": "CNY",
    "city": "novosibirsk",
    "branches": [
        {
            "bank": "Банк «Открытие»",
            "address": "630102, г. Новосибирск, ул. Кирова, дом. 44",
            "subway": "м. Октябрьская",
            "currency": "CNY",
            "buy": 9.61,
            "sell": 11.64,
            "updated": "2023-01-24T16:54:00+03:00"
        }
    ]
}
```
See [main.go](./_example/main.go).

## References

For more information check out the following links:

* Cash currency exchange rates by [Banki.ru](https://www.banki.ru/products/currency/cash/moskva/) (RU)
