package utils

import (
	"encoding/json"
	"os"
	"time"

	"github.com/go-telegram/bot/models"
	bolt "go.etcd.io/bbolt"
)

// User.
type User struct {
	User *models.User `json:"user"`
	Date time.Time    `json:"date"`
}

// Persist user.
func Persist(user *models.User) (err error) {
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
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) (err error) {
		b, err := tx.CreateBucketIfNotExists([]byte("root"))
		if err != nil {
			return
		}

		b, err = b.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return
		}

		if err = b.Put(id, u); err != nil {
			return
		}

		return
	})

	return
}
