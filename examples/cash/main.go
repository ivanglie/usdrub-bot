//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/ivanglie/usdrub-bot/pkg/bankiru-go"
)

func main() {
	client := bankiru.NewClient()
	rates, err := client.Rates(bankiru.Novosibirsk)
	if err != nil {
		panic(err)
	}
	fmt.Println(rates)
}
