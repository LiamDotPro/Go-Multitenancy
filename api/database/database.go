package database

import (
	"github.com/jinzhu/gorm"
	"github.com/wader/gormstore"
	"fmt"
	"time"
)

var Connection *gorm.DB
var Store *gormstore.Store

func init() {

	// Database Connection string
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=zeus dbname=default password=1234 sslmode=disable")

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	// Turn logging for the database on.
	db.LogMode(true)

	Connection = db

	// Now Setup store
	// @todo Add Env Variable for password.
	// Password is passed as byte key method
	Store = gormstore.NewOptions(db, gormstore.Options{
		TableName:       "sessions",
		SkipCreateTable: false,
	}, []byte("secret-hash-key"))

	// Makes quit Available
	quit := make(chan struct{})

	// Every hour remove dead sessions.
	go Store.PeriodicCleanup(1*time.Hour, quit)
}
