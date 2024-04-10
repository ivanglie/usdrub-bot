# Golang client for the Central Bank of the Russian Federation currency rates API

Golang client for the [CBRF](http://www.cbr.ru/development/).

## Example

First, ensure the library is installed and up to date by running

```
go get -u github.com/ivanglie/usdrub-bot/pkg/cbr-go
```

This is a very simple app that just displays exhange rate of US dollar.

```golang
package main

import (
	"fmt"
	"time"

	"github.com/ivanglie/usdrub-bot/pkg/cbr-go"
)

func main() {
	client := cbr.NewClient()
	rate, err := client.GetRate("USD", time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(rate)
}
```

Console output:

```
76.8207
```

See [main.go](../../examples/cbr/main.go).

## References

For more information check out the following links:

* [CBRF API](http://www.cbr.ru/development/SXML/)
* [CBRF technical resources](http://www.cbr.ru/eng/development/) (EN)