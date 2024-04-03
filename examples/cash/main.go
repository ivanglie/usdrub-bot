//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/ivanglie/usdrub-bot/pkg/go-br-client"
)

func main() {
	client := br.NewClient()
	rates, err := client.Rates(br.Novosibirsk)
	if err != nil {
		panic(err)
	}
	fmt.Println(rates)
}
