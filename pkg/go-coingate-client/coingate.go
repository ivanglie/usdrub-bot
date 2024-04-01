package coingate

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

const baseURL = "https://api.coingate.com/v2/rates/merchant"

// Debug mode
// If this variable is set to true, debug mode activated for the package
var Debug = false

// CoinGate API error responses
// Response example:
//
//	{
//	  "message": "Not found App by Access-Key",
//	  "reason": "BadCredentials"
//	}
//
// See https://developer.coingate.com/docs/common-errors
type Error struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

// Current exchange rate for any two currencies, fiat or crypto.
// This endpoint is public, authentication is not required.
// Arguments are ISO Symbol. Example: EUR, USD, BTC, ETH, etc.
// See https://developer.coingate.com/docs/get-rate
func getRate(from, to string, fetch FetchFunction) (float64, error) {
	if Debug {
		log.Printf("Fetching the currency rate for %s\n", to)
	}

	var res float64 = 0
	url := fmt.Sprintf("%s/%s/%s", baseURL, from, to)
	resp, err := fetch(url)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if resp.StatusCode != 200 {
		var error Error
		err = json.Unmarshal(b, &error)
		if err != nil {
			return res, err
		}

		return res, fmt.Errorf("service error (message: %s, reason: %s)", error.Message, error.Reason)
	}

	s := string(b)
	if res, err = strconv.ParseFloat(s, 64); err != nil {
		return res, err
	}

	return res, nil
}
