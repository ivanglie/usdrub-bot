package bestchange

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClient_Rate(t *testing.T) {
	c := NewClient()
	c.buildURL = func() string {
		dir, _ := os.Getwd()
		return "file:" + filepath.Join(dir, "/test/bestchangecom")
	}

	got, err := c.Rate()
	if err != nil {
		t.Error(err)
	}

	want := 96.414084
	if got != want {
		t.Errorf("Avg rate = %v, want %v", got, want)
	}
}

func Test_buildURL(t *testing.T) {
	buildURL := func() string {
		return baseURL
	}

	want := (NewClient()).buildURL()

	if got := buildURL(); got != want {
		t.Errorf("URL.build() = %v, want %v", got, want)
	}
}
