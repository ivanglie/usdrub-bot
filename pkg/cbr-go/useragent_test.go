package cbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_randomUserAgent(t *testing.T) {
	userAgent := randomUserAgent()
	assert.Contains(t, userAgent, "Mozilla/")
}
