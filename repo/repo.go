package repo

import (
	"github.com/boltdb/bolt"
	"time"
)

var db *bolt.DB
var usersBucket = []byte("users")
var bookmarksBucket = []byte("bookmarks")
var sessionsBucket = []byte("sessions")

func Init(dbPath string) (*bolt.DB, error) {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(usersBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(sessionsBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(bookmarksBucket)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
