package bestchange

import (
	"encoding/json"
	"fmt"
)

// Currency type.
type Currency string

// City type.
type City string

// Average exchange rate.
type Rate struct {
	Currency Currency `json:"currency"`
	City     City     `json:"city"`
	Value    float64  `json:"value"`
}

// String representation of average exchange rate.
func (r *Rate) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(b)
}
