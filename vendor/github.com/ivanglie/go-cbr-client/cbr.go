package cbr

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

const (
	baseURL    = "http://www.cbr.ru/scripts/XML_daily_eng.asp"
	dateFormat = "02/01/2006"
)

// Debug mode.
// If this variable is set to true, debug mode activated for the package.
var Debug = false

// Currency is a currency item.
type Currency struct {
	ID       string `xml:"ID,attr"`
	NumCode  uint   `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nom      uint   `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

// Result is a result representation.
type Result struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Date       string     `xml:"Date,attr"`
	Currencies []Currency `xml:"Valute"`
}

type httpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a rates service client.
type Client struct {
	httpClient httpClientInterface
}

// NewClient creates a new rates service instance.
func NewClient() *Client {
	return &Client{httpClient: &http.Client{}}
}

// GetRate returns a currency rate for a given currency and date.
func (s *Client) GetRate(currency string, t time.Time) (float64, error) {
	rate, err := s.rate(currency, t, s.httpClient)
	if err != nil {
		return 0, err
	}
	return rate, nil
}

func (s *Client) rate(currency string, t time.Time, hc httpClientInterface) (float64, error) {
	if Debug {
		log.Printf("Fetching the currency rate for %s at %v\n", currency, t.Format("02.01.2006"))
	}

	var result Result
	if err := s.currencies(&result, t); err != nil {
		return 0, err
	}

	for _, v := range result.Currencies {
		if v.CharCode == currency {
			return currencyRateValue(v)
		}
	}

	return 0, fmt.Errorf("unknown currency: %s", currency)
}

func (s *Client) currencies(v *Result, t time.Time) error {
	url := baseURL + "?date_req=" + t.Format(dateFormat)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", randomUserAgent())

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	err = decoder.Decode(&v)
	if err != nil {
		return err
	}

	return nil
}

func currencyRateValue(cur Currency) (float64, error) {
	var res float64 = 0
	properFormattedValue := strings.Replace(cur.Value, ",", ".", -1)
	res, err := strconv.ParseFloat(properFormattedValue, 64)
	if err != nil {
		return res, err
	}

	return res / float64(cur.Nom), nil
}
