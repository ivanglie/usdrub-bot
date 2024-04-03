package exchange

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_rate_Update(t *testing.T) {
	r := Get()
	r.Value(Forex).f = func() (float64, error) { return 50.0, nil }

	r.Update()
	assert.Equal(t, 50.0, r.Value(Forex).value)

	// Error
	r.Value(Forex).f = func() (float64, error) { return 51.0, errors.New("error") }

	r.Update()
	assert.Equal(t, 50.0, r.Value(Forex).value)
}

func Test_rate_String(t *testing.T) {
	r := Get()
	r.Value(Forex).value = 50.0
	r.Value(MOEX).value = 51.0
	r.Value(CBRF).value = 52.0

	t.Log(r)

	assert.Equal(t, "50.00 RUB by Forex\n51.00 RUB by Moscow Exchange\n52.00 RUB by Russian Central Bank\n", r.String())
}

func Test_rates_Value(t *testing.T) {
	r := Get()
	r.Value(Forex).value = 50.0

	assert.Equal(t, 50.0, r.Value(Forex).value)

	// Error
	assert.Nil(t, r.Value("Phorex"))
}
