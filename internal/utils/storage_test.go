package utils

import (
	"testing"

	"github.com/go-telegram/bot/models"
)

func TestPersist(t *testing.T) {
	usr := &models.User{ID: 1, FirstName: "Test", LastName: "Test"}
	if err := Persist(usr); err != nil {
		t.Errorf("Persist() error = %v", err)
	}
}
