package users

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User
type User struct {
	gorm.Model
	Email         string
	Password      string `json:",omitempty"`
	AccountType   int
	CompanyId     int
	FirstName     string
	LastName      string
	PhoneNumber   string
	RecoveryEmail string
	Contractor    bool   `gorm:"sql:'not null' default:'false'"`
}

/**
Crud
 */

// Creates a standard user in the database.
// Returns the inserted user id
func createUser(email string, password string, accountType int) (uint, error) {

	// Slice for found users.
	var foundUsers []User

	if err := database.Connection.Select("email").Where("email = ?", email).Find(&foundUsers).Error; err != nil {
		return 0, err
	}

	// If duplicate email address has been found return.
	if len(foundUsers) > 0 {
		return 0, errors.New("A user with that email address already exists")
	}

	// Hash the password so it's not clear text.
	// Run outside of the if statements so we can grab the result outside of local scope.
	hash, hashErr := hashPassword([]byte(password))

	if hashErr != nil {
		return 0, hashErr
	}

	var user = User{Email: email, Password: hash, AccountType: accountType}

	// Run create
	if err := database.Connection.Create(&user).Error; err != nil {
		// Error Handler
		return 0, err
	}

	// Return newly created user ID
	return user.ID, nil
}

// Logs a user in.
func loginUser(email string, password string) (uint, bool, error) {

	// Create local state user
	var user User

	// Find the user by email, return error if input is malformed.
	if err := database.Connection.First(&user, "email = ?", email).Error; err != nil {
		return 0, false, err
	}

	// Now we've found a user send off the hashed password and sent password for decoding.
	if result := CheckPasswordHash(password, user.Password); result != true {
		// Passwords do not match
		return 0, false, errors.New("passwords did not match")
	}

	// Checks have bee passed return true
	return user.ID, true, nil
}

// Updates a user in the database.
// A separate method is called when updating a company id
func updateUser(id uint, email string, accountType int, firstName string, lastName string, phoneNumber string, recoveryEmail string, contractor bool) (string, error) {

	var user User

	// Update the basic user information, anything that was set as nil will not be changed.
	if err := database.Connection.Model(&user).Where("id = ?", id).Updates(User{
		Email:         email,
		AccountType:   accountType,
		FirstName:     firstName,
		LastName:      lastName,
		PhoneNumber:   phoneNumber,
		RecoveryEmail: recoveryEmail,
		Contractor:    contractor,
	}).Error; err != nil {
		return "", err
	}

	return "User Information Successfully Updated.", nil

}

// Used to directly update a users companyID, this is for company assignment.
func updateUserCompany(id uint, companyId int) (string, error) {
	var user User

	if err := database.Connection.Model(&user).Where("id = ?", id).Updates(User{CompanyId: companyId}).Error; err != nil {
		return "", err
	}

	return "User Company Information Successfully Updated.", nil
}

// Deletes a user in the database.
func deleteUser(id uint) (string, error) {
	var user User

	if err := database.Connection.Where("id = ?", id).Delete(&user).Error; err != nil {
		return "An error occurred when trying to delete the user", err
	}

	return "The user has been successfully deleted", nil
}

// Get a specific user from the database.
func getUser(id uint) (*User, error) {

	var user User

	if err := database.Connection.Select("id, created_at, updated_at, deleted_at, email, account_type, company_id, first_name, last_name").Where("id = ? ", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Get all users from the database.
// @todo convert to pagination.
func getUserList() (*[]User, error) {
	var user []User

	if err := database.Connection.Select("id, created_at, updated_at, deleted_at, email, account_type, company_id, first_name, last_name").Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Gets all users who are part of a certain company.
// @todo convert to pagination.
func getUserListByCompanyId(companyId uint) (*[]User, error) {
	var user []User

	if err := database.Connection.Select("id, created_at, updated_at, deleted_at, email, account_type, company_id, first_name, last_name").Where("company_id = ?", companyId).Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Hashes a password
// Running 4 rounds to comply with bcrypt recommendations for standard user.
func hashPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 4)
	return string(hash), err
}

// Hashes a password for a super user
// Running 8 rounds to comply with bcrypt recommendations for a super user.
func hashPasswordAdmin(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 8)

	return string(hash), err
}

// Check a users password from a hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
