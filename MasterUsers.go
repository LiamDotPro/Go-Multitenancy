package main

import (
	"errors"
	"github.com/LiamDotPro/Go-Multitenancy/helpers"
	"github.com/LiamDotPro/Go-Multitenancy/tenants"
	"github.com/jinzhu/gorm"
	"strings"
)

type MasterUsers struct {
	gorm.Model
	Email         string
	Password      string `json:",omitempty"`
	AccountType   int
	FirstName     string
	LastName      string
	PhoneNumber   string
	RecoveryEmail string
}

// Creates a standard user in the database.
// Returns the inserted user id
func createMasterUser(email string, password string, accountType int) (uint, error) {

	// Slice for found users.
	var foundUsers []User

	if err := Connection.Select("email").Where("email = ?", email).Find(&foundUsers).Error; err != nil {
		return 0, err
	}

	// If duplicate email address has been found return.
	if len(foundUsers) > 0 {
		return 0, errors.New("A user with that email address already exists")
	}

	// Hash the password so it's not clear text.
	// Run outside of the if statements so we can grab the result outside of local scope.
	hash, hashErr := helpers.HashPassword([]byte(password))

	if hashErr != nil {
		return 0, hashErr
	}

	var user = User{Email: email, Password: hash, AccountType: accountType}

	// Run create
	if err := Connection.Create(&user).Error; err != nil {
		// Error Handler
		return 0, err
	}

	// Return newly created user ID
	return user.ID, nil
}

// Logs a user in.
func loginMasterUser(email string, password string) (uint, bool, error) {

	// Create local state user
	var user User

	// Find the user by email, return error if input is malformed.
	if err := Connection.First(&user, "email = ?", email).Error; err != nil {
		return 0, false, err
	}

	// Now we've found a user send off the hashed password and sent password for decoding.
	if result := helpers.CheckPasswordHash(password, user.Password); result != true {
		// Passwords do not match
		return 0, false, errors.New("passwords did not match")
	}

	// Checks have bee passed return true
	return user.ID, true, nil
}

// Updates a user in the database.
// A separate method is called when updating a company id
func updateMasterUser(id uint, email string, accountType int, firstName string, lastName string, phoneNumber string, recoveryEmail string) (string, error) {

	var user User

	// Update the basic user information, anything that was set as nil will not be changed.
	if err := Connection.Model(&user).Where("id = ?", id).Updates(User{
		Email:         email,
		AccountType:   accountType,
		FirstName:     firstName,
		LastName:      lastName,
		PhoneNumber:   phoneNumber,
		RecoveryEmail: recoveryEmail,
	}).Error; err != nil {
		return "", err
	}

	return "User Information Successfully Updated.", nil

}

// Deletes a user in the database.
func deleteMasterUser(id uint) (string, error) {
	var user User

	if err := Connection.Where("id = ?", id).Delete(&user).Error; err != nil {
		return "An error occurred when trying to delete the user", err
	}

	return "The user has been successfully deleted", nil
}

// Create's a tenant using a domain identifier
func createNewTenant(subDomainIdentifier string) (msg string, err error) {

	// Create new database to hold client.
	if err := Connection.Exec("CREATE DATABASE " + strings.ToLower(subDomainIdentifier) + " OWNER admin").Error; err != nil {
		return "error making the database", err
	}

	var connectionInfo = tenants.TenantConnectionInformation{TenantSubDomainIdentifier: subDomainIdentifier, ConnectionString: "host=localhost port=5432 user=admin dbname=" + subDomainIdentifier + " password=1234 sslmode=disable"}

	if err := Connection.Create(&connectionInfo).Error; err != nil {
		return "error inserting the new database record", err
	}

	tenConn, tenConErr := connectionInfo.GetConnection()

	if tenConErr != nil {
		return "error creating the connection using connection method", err
	}

	if migrateErr := migrateTenantTables(tenConn); migrateErr != nil {
		return "error attempting to migrate the existing tables to new database", migrateErr
	}

	return "New Tenant has been successfully made", nil
}

// Get a specific user from the database.
func getMasterUser(id uint) (*User, error) {

	var user User

	if err := Connection.Select("id, created_at, updated_at, deleted_at, email, account_type, company_id, first_name, last_name").Where("id = ? ", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
