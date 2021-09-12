package repo

import (
	"encoding/json"
	"github.com/saarmornel/reading-list/misc"
	//	log "github.com/sirupsen/logrus"
)

type Bookmark struct {
	Url     string   `json:"url,omitempty"`
	Read    bool     `json:"read,omitempty"`
	Private bool     `json:"private,omitempty"`
	Title   string   `json:"title,omitempty"`
	Icon    string   `json:"icon,omitempty"`
	Comment string   `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

func GetBookmarks(username string) ([]Bookmark, error) {
	var bookmarks []Bookmark
	tx, err := db.Begin(true)
	if err != nil {
		return bookmarks, err
	}
	defer tx.Rollback()

	root := tx.Bucket(bookmarksBucket)
	bkt := root.Bucket([]byte(username))
	bkt.ForEach(func(k, v []byte) error {
		var b Bookmark
		json.Unmarshal(v, &b)
		bookmarks = append(bookmarks, b)
		return nil
	})

	return bookmarks, nil
}

func CreateBookmark(username string, b *Bookmark) (*Bookmark, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	root := tx.Bucket(bookmarksBucket)
	bkt, err := root.CreateBucketIfNotExists([]byte(username))
	if err != nil {
		return nil, err
	}
	if buf, err := json.Marshal(b); err != nil {
		return nil, err
	} else if err := bkt.Put([]byte(b.Url), buf); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return b, nil
}

func DeleteBookmark(username string, url string) (*misc.Status, error) {
	status := &misc.Status{}
	tx, err := db.Begin(true)
	if err != nil {
		status.SetStatus(false)
		return status, err
	}
	defer tx.Rollback()

	root := tx.Bucket(bookmarksBucket)
	bkt := root.Bucket([]byte(username))

	if err := bkt.Delete([]byte(url)); err != nil {
		status.SetStatus(false)
		return status, err
	}

	if err := tx.Commit(); err != nil {
		status.SetStatus(false)
		return status, err
	}
	status.SetStatus(true)
	return status, nil
}
