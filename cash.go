package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Provides code safe
var cmu sync.RWMutex

// Branch
type Branch struct {
	Bank    string  `json:"bank"`
	Address string  `json:"address"`
	Subway  string  `json:"subway"`
	Code    string  `json:"code"`
	Buy     float64 `json:"buy"`
	Sell    float64 `json:"sell"`
	Updated string  `json:"updated"`
}

// Cash exchange rate
type Cash struct {
	name     string
	pattern  string
	branches []Branch
	min,
	max,
	avg float64
}

// Update rate from source
func (c *Cash) updateRate(url string) error {
	cmu.Lock()
	defer cmu.Unlock()

	if len(url) == 0 {
		return fmt.Errorf("no url set")
	}

	log.Info("Fetching the cash currency rate for RUB")

	res, err := http.Get(url + "?region=moskva")
	if err != nil {
		log.Errorf("Error making http request: %s", err)
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error making read: %s", err)
		return err
	}

	if res.StatusCode != 200 {
		log.Errorf("Server error (status code: %d)", res.StatusCode)
		c.min = 0.0
		c.max = 0.0
		c.avg = 0.0
		return fmt.Errorf("server error (%+v): %v", res.StatusCode, res)
	}

	var b []Branch
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Errorf("Error making unmarshal: %s", err)
		return err
	}

	c.branches = b

	c.min = c.getMin()
	c.max = c.getMax()
	c.avg = c.getAvg()

	return nil
}

// Get formated cash exchange rates from Cash struct
func (c *Cash) getRatef() string {
	cmu.RLock()
	defer cmu.RUnlock()
	res := fmt.Sprintf(c.pattern, c.min, c.max, c.avg)
	if c.min <= 0.0 || c.max <= 0.0 || c.avg <= 0.0 {
		res = fmt.Sprintf("%s error: wrong value of min=%.2f, max=%.2f, avg=%.2f", c.name, c.min, c.max, c.avg)
	}
	return res
}

// Min
func (c *Cash) getMin() float64 {
	min := c.branches[0].Sell
	for _, v := range c.branches {
		if v.Sell < min {
			min = v.Sell
		}
	}
	return min
}

// Max
func (c *Cash) getMax() float64 {
	max := c.branches[0].Sell
	for _, v := range c.branches {
		if v.Sell > max {
			max = v.Sell
		}
	}
	return max
}

// Avg
func (c *Cash) getAvg() float64 {
	total := 0.0
	for _, b := range c.branches {
		total += b.Sell
	}
	return total / float64(len(c.branches))
}
