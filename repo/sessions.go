package repo

import (
	"github.com/boltdb/bolt"
)

func CreateSession(token string, username string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sessions"))
		return b.Put([]byte(token), []byte(username))
	})
	if err != nil {
		return err
	}
	return nil
}

func GetSession(token string) (string, error) {
	var username string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("sessions"))
		username = string(b.Get([]byte(token)))
		return nil
	})
	if err != nil {
		return "", err
	}
	return username, nil
}
