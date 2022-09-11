package moex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	baseURL = "https://iss.moex.com/iss"
)

// Codes
const (
	USDRUB = "USD000UTSTOM"
	EURRUB = "EUR_RUB__TOM"
	GBPRUB = "GBPRUB_TOM"
	CNYRUB = "CNYRUB_TOM"
)

// Debug mode
// If this variable is set to true, debug mode activated for the package
var Debug = false

type Currency struct {
	values []currency
}

type currency struct {
	Charsetinfo struct {
		Name string `json:"name"`
	} `json:"charsetinfo,omitempty"`
	Securities []struct {
		Secid       string      `json:"SECID"`
		Boardid     string      `json:"BOARDID"`
		Shortname   string      `json:"SHORTNAME"`
		Lotsize     int         `json:"LOTSIZE"`
		Settledate  string      `json:"SETTLEDATE"`
		Decimals    int         `json:"DECIMALS"`
		Facevalue   int         `json:"FACEVALUE"`
		Marketcode  string      `json:"MARKETCODE"`
		Minstep     float64     `json:"MINSTEP"`
		Prevdate    string      `json:"PREVDATE"`
		Secname     string      `json:"SECNAME"`
		Remarks     interface{} `json:"REMARKS"`
		Status      string      `json:"STATUS"`
		Faceunit    string      `json:"FACEUNIT"`
		Prevprice   float64     `json:"PREVPRICE"`
		Prevwaprice float64     `json:"PREVWAPRICE"`
		Currencyid  string      `json:"CURRENCYID"`
		Latname     string      `json:"LATNAME"`
		Lotdivider  int         `json:"LOTDIVIDER"`
	} `json:"securities,omitempty"`
	Marketdata []struct {
		Highbid               interface{} `json:"HIGHBID"`
		Biddepth              interface{} `json:"BIDDEPTH"`
		Lowoffer              interface{} `json:"LOWOFFER"`
		Offerdepth            interface{} `json:"OFFERDEPTH"`
		Spread                float64     `json:"SPREAD"`
		High                  float64     `json:"HIGH"`
		Low                   float64     `json:"LOW"`
		Open                  float64     `json:"OPEN"`
		Last                  float64     `json:"LAST"`
		Lastcngtolastwaprice  float64     `json:"LASTCNGTOLASTWAPRICE"`
		Valtoday              float64     `json:"VALTODAY"`
		Voltoday              float64     `json:"VOLTODAY"`
		ValtodayUsd           float64     `json:"VALTODAY_USD"`
		Waprice               float64     `json:"WAPRICE"`
		Waptoprevwaprice      float64     `json:"WAPTOPREVWAPRICE"`
		Closeprice            interface{} `json:"CLOSEPRICE"`
		Numtrades             int         `json:"NUMTRADES"`
		Tradingstatus         string      `json:"TRADINGSTATUS"`
		Updatetime            string      `json:"UPDATETIME"`
		Boardid               string      `json:"BOARDID"`
		Secid                 string      `json:"SECID"`
		Waptoprevwapriceprcnt float64     `json:"WAPTOPREVWAPRICEPRCNT"`
		Bid                   interface{} `json:"BID"`
		Biddeptht             interface{} `json:"BIDDEPTHT"`
		Numbids               interface{} `json:"NUMBIDS"`
		Offer                 interface{} `json:"OFFER"`
		Offerdeptht           interface{} `json:"OFFERDEPTHT"`
		Numoffers             interface{} `json:"NUMOFFERS"`
		Change                float64     `json:"CHANGE"`
		Lastchangeprcnt       float64     `json:"LASTCHANGEPRCNT"`
		Value                 float64     `json:"VALUE"`
		ValueUsd              float64     `json:"VALUE_USD"`
		Seqnum                int64       `json:"SEQNUM"`
		Qty                   int         `json:"QTY"`
		Time                  string      `json:"TIME"`
		Priceminusprevwaprice float64     `json:"PRICEMINUSPREVWAPRICE"`
		Lastchange            float64     `json:"LASTCHANGE"`
		Lasttoprevprice       float64     `json:"LASTTOPREVPRICE"`
		ValtodayRur           int64       `json:"VALTODAY_RUR"`
		Systime               string      `json:"SYSTIME"`
		Marketprice           float64     `json:"MARKETPRICE"`
		Marketpricetoday      float64     `json:"MARKETPRICETODAY"`
		Marketprice2          interface{} `json:"MARKETPRICE2"`
		Admittedquote         interface{} `json:"ADMITTEDQUOTE"`
		Lopenprice            float64     `json:"LOPENPRICE"`
	} `json:"marketdata,omitempty"`
}

// Current exchange rate for two currencies.
// This endpoint is public, authentication is not required.
// Example: EURRUB, USDRUB, etc.
// See https://iss.moex.com/iss/reference/
func getRate(code string, fetch fetchFunction) (float64, error) {
	if Debug {
		log.Printf("Fetching the currency rate for %s\n", code)
	}

	var res float64 = 0
	url := fmt.Sprintf("%s%s%s", baseURL,
		"/engines/currency/markets/selt/securities.json?iss.only=securities,marketdata&lang=en&iss.meta=off&iss.json=extended",
		"&securities=CETS:"+code)
	resp, err := fetch(url)
	if err != nil {
		return res, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	var c *Currency = &Currency{}
	err = json.Unmarshal(body, &c.values)
	if err != nil {
		return res, err
	}
	if c == (&Currency{}) {
		return res, fmt.Errorf("MOEX error: c is zero")
	}
	if len(c.values) < 2 {
		return res, fmt.Errorf("MOEX error: length of c.values less than 2")
	}
	val := c.values[1]
	if val.Marketdata == nil {
		return res, fmt.Errorf("MOEX error: val.Marketdata is zero")
	}
	md := val.Marketdata
	if len(md) == 0 {
		return 0, fmt.Errorf("MOEX error: length of md equals 0")
	}
	res = md[0].Last
	return res, nil
}
