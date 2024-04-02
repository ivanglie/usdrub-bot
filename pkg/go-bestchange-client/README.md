# Golang client for bestchange.com

[![Go Reference](https://pkg.go.dev/badge/github.com/ivanglie/go-bestchange-client.svg)](https://pkg.go.dev/github.com/ivanglie/go-bestchange-client)
[![Test](https://github.com/ivanglie/go-bestchange-client/actions/workflows/test.yml/badge.svg)](https://github.com/ivanglie/go-bestchange-client/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/ivanglie/go-bestchange-client/branch/master/graph/badge.svg?token=8lRyze5RSQ)](https://codecov.io/gh/ivanglie/go-bestchange-client)

Golang client that provides latest exchange rate.

## Example

First, ensure the library is installed and up to date by running

```
go get -u github.com/ivanglie/go-bestchange-client
```

This is a very simple app that just displays Cash RUB to Tether TRC20 (USDT) exhange rate in Moscow.

```golang
package main

import (
	"fmt"

	bestchange "github.com/ivanglie/go-bestchange-client"
)

func main() {
	client := bestchange.NewClient()
	rate, err := client.Rate(bestchange.Moscow)
	if err != nil {
		panic(err)
	}
	fmt.Println(rate)
}
```

Console output:

```json
{"currency":"USDT","city":"msk","value":95.841441}
```
See [main.go](./_example/main.go).

## References

For more information check out the following links:

* Exchange Cash RUB to Tether TRC20 (USDT) in Moscow [bestchange.com](https://www.bestchange.com/cash-ruble-to-tether-trc20-in-msk.html)
