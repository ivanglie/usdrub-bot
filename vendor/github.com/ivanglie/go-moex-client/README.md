# Golang client for Moscow Exchange ISS API

[![Go Reference](https://pkg.go.dev/badge/github.com/ivanglie/go-moex-client.svg)](https://pkg.go.dev/github.com/ivanglie/go-moex-client)
[![Test](https://github.com/ivanglie/go-moex-client/actions/workflows/test.yml/badge.svg)](https://github.com/ivanglie/go-moex-client/actions/workflows/test.yml)

Golang client for the [Moscow Exchange ISS API](https://www.moex.com/a2193).

## Example

First, ensure the library is installed and up to date by running ```go get -u github.com/ivanglie/go-moex-client```.

This is a very simple app that just displays Chinese Yuan Renminbi to Russian Ruble conversion.

```golang
package main

import (
	"fmt"

	moex "github.com/ivanglie/go-moex-client"
)

func main() {
	client := moex.NewClient()
	rate, err := client.GetRate(moex.CNYRUB)
	if err != nil {
		panic(err)
	}
	fmt.Println(rate)
}
```
See [main.go](./_example/main.go).

## References

For more information check out the following links:

* [MOEX ISS API](https://www.moex.com/a2193)
* [MOEX ISS reference](https://iss.moex.com/iss/reference/)