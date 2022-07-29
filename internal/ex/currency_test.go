package ex

import (
	"math/rand"
	"testing"
	"time"
)

func TestDataRace(t *testing.T) {
	c := New("%.2f", func() (float64, error) { return 100 * rand.Float64(), nil })
	go func() {
		for {
			c.Update()
		}
	}()

	for i := 0; i < 10; i++ {
		c.Format()
		time.Sleep(100 * time.Millisecond)
	}
}
