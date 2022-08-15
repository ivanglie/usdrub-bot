package ex

import (
	"math/rand"
	"testing"
	"time"
)

func Test_dataRace(t *testing.T) {
	c := New(func() (float64, error) { return 100 * rand.Float64(), nil })
	go func() {
		for {
			c.Update(nil)
		}
	}()

	for i := 0; i < 10; i++ {
		c.Rate()
		time.Sleep(100 * time.Millisecond)
	}
}
