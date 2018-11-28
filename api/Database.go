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
var TenantInformation []TenantConnectionInformation

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

	msg, err := createNewTenant("Happy Feet")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(msg)

	// Get all connection information from the database
	getTenantDataFromDatabase()

	// Makes quit Available
	quit := make(chan struct{})

	// Every hour remove dead sessions.
	go Store.PeriodicCleanup(1*time.Hour, quit)
}

// Create's a new database for use as a sub client.
func createNewTenant(subDomainIdentifier string) (msg string, err error) {

	// Create new database to hold client.
	if err := Connection.Exec("CREATE DATABASE " + subDomainIdentifier + " OWNER admin").Error; err != nil {
		return "", err
	}

	var connectionInfo = TenantConnectionInformation{TenantSubDomainIdentifier: subDomainIdentifier, ConnectionString: "host=localhost port=5432 user=admin dbname=" + subDomainIdentifier + " password=1234 sslmode=disable"}

	if err := Connection.Create(&connectionInfo).Error; err != nil {
		return "", err
	}

	tenConn, tenConErr := connectionInfo.getConnection()

	if tenConErr != nil {
		return "", err
	}

	if err := migrateTenantTables(tenConn); err != nil {
		return "", err
	}

	// Add the newly created tenant id back onto the tenant object

	// Add the new tenant info to the collection
	TenantInformation = append(TenantInformation, connectionInfo)

	return "New Tenant has been successfully made", nil
}

// Gets all of the current client information from the master database and loads the id's
// Into the connection information slice, then
func getTenantDataFromDatabase() {
	Connection.Find(&TenantInformation)
	fmt.Println(&TenantInformation)
}

/**
This method uses the base tenant connection set out within init.
 */
func migrateMasterTenantDatabase() error {

	if err := Connection.AutoMigrate(&TenantConnectionInformation{}).Error; err != nil {
		return err
	}

	return nil

}

// Attempts to migrate tables using database connection
func migrateTenantTables(connection *gorm.DB) error {
	if err := connection.AutoMigrate(&User{}).Error; err != nil {
		return err
	}

	return nil
}
