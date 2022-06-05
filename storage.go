package main

import (
	"encoding/json"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	bolt "go.etcd.io/bbolt"
)

// User data
type User struct {
	User *tgbotapi.User `json:"user"`
	Date time.Time      `json:"date"`
}

// Persist data
func persist(user *tgbotapi.User) error {
	id, err := json.Marshal(user.ID)
	if err != nil {
		return err
	}
	usr := User{User: user, Date: time.Now().Local()}
	u, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	db, err := bolt.Open("users.db", 0600, nil)
	if err != nil {
		log.Error(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		r, err := tx.CreateBucketIfNotExists([]byte("root"))
		if err != nil {
			return err
		}
		b, err := r.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		err = b.Put([]byte(id), []byte(u))
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return err
	}
	return nil
}
