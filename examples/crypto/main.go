//go:build ignore
// +build ignore

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
