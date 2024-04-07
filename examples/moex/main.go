//go:build ignore
// +build ignore

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
