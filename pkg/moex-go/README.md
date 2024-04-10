# Golang client for Moscow Exchange ISS API

Golang client for the [Moscow Exchange](https://www.moex.com/a2920).

## Example

First, ensure the library is installed and up to date by running 

```
go get -u github.com/ivanglie/usdrub-bot/pkg/moex-go
```

This is a very simple app that just displays Chinese Yuan Renminbi to Russian Ruble conversion.

```golang
package main

import (
	"fmt"

	moex "github.com/ivanglie/usdrub-bot/pkg/moex-go"
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
See [main.go](../../examples/moex/main.go).

## References

For more information check out the following links:

* MOEX ISS API [en](https://www.moex.com/a2920), [ru](https://www.moex.com/a2193)
* [MOEX ISS reference](https://iss.moex.com/iss/reference/)