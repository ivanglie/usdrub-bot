//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/ivanglie/usdrub-bot/pkg/go-cbr-client"
)

func main() {
	client := cbr.NewClient()
	rate, err := client.GetRate("USD", time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(rate)
}
