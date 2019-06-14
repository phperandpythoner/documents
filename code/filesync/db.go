package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

//数据库
const dbFile = "file_change.db"

//数据库表
const bucket = "fileBucket"

type FileDB struct {
	db   *bolt.DB
	data []byte
	//operateType int
}

//check file does it exist
func dbExists(dbfile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

//create db
func createDB(dbfile string) bool {
	_, err := os.OpenFile(dbFile, os.O_RDWR|os.O_CREATE, 0777)
	if err == nil {
		return true
	}

	return false
}

func (db *FileDB) initDB() *bolt.DB {
	if dbExists(dbFile) == false {
		if createDB(dbFile) == false {
			fmt.Printf("create database %s fial", dbFile)
			os.Exit(1)
		}
	}

	boltDB, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	//defer boltDB.Close()

	return boltDB
}

func (db *FileDB) DB() *FileDB {
	boltDB := db.initDB()
	err := boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))

		return err
	})
	if err != nil {
		log.Panic(err)
	}

	return &FileDB{
		db:   boltDB,
		data: db.data,
	}
}

func (db *FileDB) Create() {
	// if dbExists(dbFile) == false {
	// 	fmt.Printf("database %s not exists.", dbFile)
	// 	os.Exit(1)
	// }

	// boltDB, err := bolt.Open(dbFile, 0600, nil)
	// if err != nil {
	// 	log.Panic(err)
	// }

	boltDB := db.db

	err := boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		err := b.Put(db.data, []byte("cmd"))

		return err
	})
	if err != nil {
		log.Panic(err)
	}
	boltDB.Close()
}

func (db *FileDB) Viwes() {
	// if dbExists(dbFile) == false {
	// 	fmt.Printf("database %s not exists.", dbFile)
	// 	os.Exit(1)
	// }
	// boltDB, err := bolt.Open(dbFile, 0600, nil)
	// if err != nil {
	// 	log.Panic(err)
	// }

	boltDB := db.db

	err := boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		//b.Cursor() 迭代所有数据
		c := b.Cursor()

		//遍历所有数据
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key:%s, value:%s\n", k, v)
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	boltDB.Close()
}
