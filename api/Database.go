package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/wader/gormstore"
	"strings"
	"time"
)

var Connection *gorm.DB
var Store *gormstore.Store
var TenantInformation []TenantConnectionInformation
var TenantMap map[string]TenantConnectionInformation

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
	if err := Connection.Exec("CREATE DATABASE " + strings.ToLower(subDomainIdentifier) + " OWNER admin").Error; err != nil {
		return "error making the database", err
	}

	var connectionInfo = TenantConnectionInformation{TenantSubDomainIdentifier: subDomainIdentifier, ConnectionString: "host=localhost port=5432 user=admin dbname=" + subDomainIdentifier + " password=1234 sslmode=disable"}

	if err := Connection.Create(&connectionInfo).Error; err != nil {
		return "error inserting the new database record", err
	}

	tenConn, tenConErr := connectionInfo.getConnection()

	if tenConErr != nil {
		return "error creating the connection using connection method", err
	}

	if migrateErr := migrateTenantTables(tenConn); migrateErr != nil {
		return "error attempting to migrate the existing tables to new database", migrateErr
	}

	// Add the newly created tenant id back onto the tenant object

	// Add the new tenant info to the collection
	TenantInformation = append(TenantInformation, connectionInfo)

	return "New Tenant has been successfully made", nil
}

// Gets all of the current client information from the master database and loads the id's
// Into the connection information slice, then calls migrates on all of the profiles
func getTenantDataFromDatabase() {
	Connection.Find(&TenantInformation)

	TenantMap = make(map[string]TenantConnectionInformation)

	for _, element := range TenantInformation {

		TenantMap[element.TenantSubDomainIdentifier] = element

		conn, _ := element.getConnection()
		migrateTenantTables(conn)
	}
}

/**
This method uses the base tenant connection set out within init.
 */
func migrateMasterTenantDatabase() error {

	if err := Connection.AutoMigrate(&TenantConnectionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&TenantSubscriptionInformation{}).Error; err != nil {
		return err
	}

	if err := Connection.AutoMigrate(&TenantSubscriptionType{}).Error; err != nil {
		return err
	}

	return nil

}

// Attempts to migrate tables using database connection
func migrateTenantTables(connection *gorm.DB) error {
	fmt.Println("Attempting to migrate tables to new database.")

	if err := connection.AutoMigrate(&User{}).Error; err != nil {
		return err
	}

	return nil
}
