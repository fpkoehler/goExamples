/* Experimenting with storing file paths in the DB. 
 *
 * This version stores the full path to every file.  After running from my root directory:
 *  -rw------- 1 fred fred 43282432 Jul 22 15:05 /tmp/my.db
 *    Number of keys/value pairs: 109953
 *    Total number of buckets: 2
 *
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/boltdb/bolt"
)

//var db *bolt.DB

const dbFilePath = "/tmp/my.db"

func createBuckets() error {
	// It will be created if it doesn't exist.
	db, err := bolt.Open(dbFilePath, 0600, nil)
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

	return err
}

func visitFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return nil
	}
	if info.IsDir() {
		return nil
	}

	db, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		topBucket := tx.Bucket([]byte("MyBucket"))
		filesBucket := topBucket.Bucket([]byte("Files"))
		v := filesBucket.Get([]byte(path))
		if v != nil {
			/* path aleady exists, don't increment id */
			return nil
		}
		id, _ := filesBucket.NextSequence()
		idStr := strconv.FormatUint(id, 10)
		err := filesBucket.Put([]byte(path), []byte(idStr))
		if err != nil {
			log.Fatalln("db.Update", path, err)
		}
		return err
	})

	// continue walk even if errors
	return nil
}

func printDb() error {
	db, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, topBucket *bolt.Bucket) error {
			fmt.Println("Bucket:", string(name))
			b := topBucket.Bucket([]byte("Files"))

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf(" key=%s, value=%s\n", k, v)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func main() {
	var err error

	err = createBuckets()
	if err != nil {
		log.Fatalln("Unable to create buckets:", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("Unable to get directory/folder path:", err)
	}
	err = filepath.Walk(dir, visitFile)

	//err = printDb()
	if err != nil {
		log.Fatalln("Db print failure:", err)
	}
}
