package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

var printPath []string

func recur(buk *bolt.Bucket, level int) {
    if buk == nil {
        return
    }
    c := buk.Cursor()
    for k, v := c.First(); k != nil; k, v = c.Next() {
    	if len(printPath) == level {
    		printPath = append(printPath, string(k))
    	} else {
    		printPath[level] = string(k)
    	}
        //fmt.Printf("%s:%s\n", k, v)
        if v == nil {
            recur(buk.Bucket(k), level+1)
        } else {
        	fmt.Println(filepath.Join(printPath[:level+1]...), ":", string(v))
        }
    }
}

func main() {
	var err error

	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err = bolt.Open("/tmp/my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create MyBucket: %s\n", err)
		}
		_, err = b.CreateBucketIfNotExists([]byte("Files"))
		if err != nil {
			return fmt.Errorf("create MyBucket.Files: %s\n", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, topBucket *bolt.Bucket) error {
			fmt.Println("Bucket:", string(name))
			printPath = make([]string, 0)
			recur(topBucket, 0)
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
