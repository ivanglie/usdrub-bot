package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Cash exchange rate
type Cash struct {
	mu       sync.RWMutex
	name     string
	pattern  string
	branches []Branch
	min,
	max,
	avg float64
}

// Update rate from cash
func (c *Cash) UpdateRate(url, region string) error {
	log.Info("Fetching the cash currency rate for RUB")

	if len(url) == 0 {
		return fmt.Errorf("no url set")
	}

	if len(region) == 0 {
		region = "?region=moskva"
	} else {
		region = "?region=" + region
	}
	log.Debugln(region)

	res, err := http.Get(url + region)
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
		c.SetMin(0.0)
		c.SetMax(0.0)
		c.SetAvg(0.0)
		return fmt.Errorf("server error (%+v): %v", res.StatusCode, res)
	}

	var b []Branch
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Errorf("Error making unmarshal: %s", err)
		return err
	}

	c.SetBranches(b)

	c.SetMin(c.Min())
	c.SetMax(c.Max())
	c.SetAvg(c.Avg())

	return nil
}

// Get formated cash exchange rates
func (c *Cash) GetRatef() string {
	res := fmt.Sprintf(c.GetPattern(), c.GetMin(), c.GetMax(), c.GetAvg())
	if c.GetMin() <= 0.0 || c.GetMax() <= 0.0 || c.GetAvg() <= 0.0 {
		res = fmt.Sprintf("%s error: wrong value of min=%.2f, max=%.2f, avg=%.2f", c.GetName(),
			c.GetMin(), c.GetMax(), c.GetAvg())
	}
	return res
}

// Min
func (c *Cash) Min() float64 {
	min := c.GetBranches()[0].getSell()
	for _, b := range c.GetBranches() {
		if b.getSell() < min {
			min = b.getSell()
		}
	}
	return min
}

// Maximum
func (c *Cash) Max() float64 {
	max := c.GetBranches()[0].getSell()
	for _, b := range c.GetBranches() {
		if b.getSell() > max {
			max = b.getSell()
		}
	}
	return max
}

// Average
func (c *Cash) Avg() float64 {
	total := 0.0
	for _, b := range c.GetBranches() {
		total += b.getSell()
	}
	return total / float64(len(c.GetBranches()))
}

func (c *Cash) GetName() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.name
}
func (c *Cash) GetPattern() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pattern
}

func (c *Cash) GetBranches() []Branch {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.branches
}

func (c *Cash) GetMin() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.min
}

func (c *Cash) GetMax() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.max
}

func (c *Cash) GetAvg() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.avg
}

func (c *Cash) SetBranches(branches []Branch) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.branches = branches
}

func (c *Cash) SetMin(min float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.min = min
}

func (c *Cash) SetMax(max float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.max = max
}

func (c *Cash) SetAvg(avg float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.avg = avg
}
