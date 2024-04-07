//go:build ignore
// +build ignore

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
