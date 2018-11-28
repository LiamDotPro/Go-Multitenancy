package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/wader/gormstore"
	"time"
)

var Connection *gorm.DB
var Store *gormstore.Store
// var ClientStores []
var ConnectionInformation []TenantConnectionInformation

func startDatabaseServices() {

	// Database Connection string
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=admin dbname=master password=1234 sslmode=disable")

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	// Turn logging for the database on.
	db.LogMode(true)

	// Make Master connection available globally.
	Connection = db

	// Now Setup store - Tenant Store
	// @todo Add Env Variable for password.
	// Password is passed as byte key method
	Store = gormstore.NewOptions(db, gormstore.Options{
		TableName:       "sessions",
		SkipCreateTable: false,
	}, []byte("masterKeyPairValue"))

	// Always attempt to migrate changes to the master tenant schema
	migrateMasterTenantDatabase()

	createNewTenant()

	// Get all connection information from the database
	getTenantDataFromDatabase()

	// Makes quit Available
	quit := make(chan struct{})

	// Every hour remove dead sessions.
	go Store.PeriodicCleanup(1*time.Hour, quit)
}

// Create's a new database for use as a sub client.
func createNewTenant(subDomainIdentifier string) []error {
	if err := Connection.Exec("CREATE DATABASE " + subDomainIdentifier + " OWNER admin").GetErrors(); len(err) != 0 {
		fmt.Println(len(err))

		for _, err := range err {
			fmt.Println(err)
		}

		return err
	}
}

// Gets all of the current client information from the master database and loads the id's
// Into the connection information slice, then
func getTenantDataFromDatabase() {
	Connection.Find(&ConnectionInformation)
	fmt.Println(&ConnectionInformation)
}

/**
This method uses the base tenant connection set out within init.
 */
func migrateMasterTenantDatabase() {
	Connection.AutoMigrate(&TenantConnectionInformation{})
}

// Attempts to migrate tables using database connection
func migrateTenantTables() {
	// conn, err := ConnectionInformation[connectionIdentifier].getConnection()
}
