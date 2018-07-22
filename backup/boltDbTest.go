/* Experimenting with storing file paths in the DB. 
 *
 * This version creates child buckets for each subdirectory.  After running from my root directory:
 * -rw------- 1 fred fred 48971776 Jul 22 14:55 /tmp/my.db
 *   Number of keys/value pairs: 127625
 *   Total number of buckets: 17677
 *
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

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

	dirPath := filepath.Dir(path)
	fileName := filepath.Base(path)
	dirPathSlice := strings.Split(dirPath, string(filepath.Separator))

	db, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		topBucket := tx.Bucket([]byte("MyBucket"))
		fileBucket := topBucket.Bucket([]byte("Files"))

		/* create buckets for each dir down the path */
		createBucketSkipped := 0;
		for _, d := range dirPathSlice[1:] {
			child, err := fileBucket.CreateBucket([]byte(d))
			if err == bolt.ErrBucketExists {
				fileBucket = fileBucket.Bucket([]byte(d))
				createBucketSkipped++
			} else if err != nil {
					log.Fatalln("Error creating bucket path", dirPathSlice, d, err)
			} else {
				fileBucket = child
			}
		}

		v := fileBucket.Get([]byte(fileName))
		if v != nil {
			/* path aleady exists, don't increment id */
			return nil
		}
		id, _ := fileBucket.NextSequence()
		idStr := strconv.FormatUint(id, 10)
		err := fileBucket.Put([]byte(fileName), []byte(idStr))
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
