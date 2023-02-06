package utils

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestPersist(t *testing.T) {
	usr := &tgbotapi.User{ID: 1, FirstName: "Test", LastName: "Test"}
	if err := Persist(usr); err != nil {
		t.Errorf("Persist() error = %v", err)
	}
}
