package moex

import (
	"testing"
	"time"
)

func Test_dataRace(t *testing.T) {
	c := New()
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
