//go:build ignore
// +build ignore

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
