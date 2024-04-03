package bestchange

import (
	"encoding/json"
	"fmt"
)

// Average exchange rate.
type Rate struct {
	Value float64 `json:"value"`
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
