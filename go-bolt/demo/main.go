package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

//Post ...
type Post struct {
	Create  time.Time
	Title   string
	Content string
}

func main() {
	var err error
	db, err = bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("posts"))
		if err != nil {
			log.Fatal(err)
			return err
		}
		return b.Put([]byte("2015-01-01"), []byte("My New Year post"))
	})
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("posts"))
		v := b.Get([]byte("2015-01-01"))
		fmt.Printf("%s\n", v)
		return nil
	})
}

func testPost() {
	post := &Post{
		Create:  time.Now(),
		Title:   "my first post",
		Content: "Hello, this is my first post.",
	}
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("posts"))
		if err != nil {
			log.Fatal(err)
			return err
		}
		encoded, err := json.Marshal(post)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return b.Put([]byte(post.Create.Format(time.RFC3339)), encoded)
	})
}
