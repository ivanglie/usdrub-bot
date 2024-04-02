//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/ivanglie/usdrub-bot/pkg/go-bestchange-client"
)

func main() {
	client := bestchange.NewClient()
	rate, err := client.Rate(bestchange.Moscow)
	if err != nil {
		panic(err)
	}
	fmt.Println(rate)
}
