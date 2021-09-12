package repo

import (
	"encoding/json"
	"github.com/boltdb/bolt"
)

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
}

func GetUser(username string) (*User, error) {
	var user User
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket)
		result := b.Get([]byte(username))

		if result == nil {
			return nil
		}
		err := json.Unmarshal(result, &user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(u *User) (*User, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bookmarksBucket)
		_, err := b.CreateBucket([]byte(u.Username))
		if err != nil {
			return err
		}

		b = tx.Bucket(usersBucket)
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		return b.Put([]byte(u.Username), buf)
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}
