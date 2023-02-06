package utils

import (
	"encoding/json"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	bolt "go.etcd.io/bbolt"
)

// User data.
type User struct {
	User *tgbotapi.User `json:"user"`
	Date time.Time      `json:"date"`
}

// Persist data.
func Persist(user *tgbotapi.User) (err error) {
	id, err := json.Marshal(user.ID)
	if err != nil {
		return
	}

	usr := User{User: user, Date: time.Now().Local()}
	u, err := json.Marshal(usr)
	if err != nil {
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		return
	}

	db, err := bolt.Open(wd+"/users.db", 0600, nil)
	if err != nil {
		return
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	err = db.Update(func(tx *bolt.Tx) (err error) {
		b, err := tx.CreateBucketIfNotExists([]byte("root"))
		if err != nil {
			return
		}

		b, err = b.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return
		}

		err = b.Put(id, u)
		if err != nil {
			return
		}

		return
	})

	return
}
