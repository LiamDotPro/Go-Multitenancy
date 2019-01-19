package main

import (
	"fmt"
	"github.com/LiamDotPro/Go-Multitenancy/helpers"
	"github.com/LiamDotPro/Go-Multitenancy/middleware"
	"github.com/LiamDotPro/Go-Multitenancy/params"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

// Init
func setupUsersRoutes(router *gin.Engine) {

	users := router.Group("/api/users")

	// Turn on the need for tenancy finding.
	users.Use(middleware.FindTenancy(Connection))

	// POST
	users.POST("create", HandleCreateUser)
	users.POST("login", HandleLogin)
	users.POST("updateUserDetails", HandleUpdateUserDetails)
	users.POST("testPoster", HandleTestPoster)

	// GET
	users.GET("getUserById", HandleGetUserById)
	users.GET("getCurrentUser", HandleGetCurrentUser)
	users.GET("testGetter", HandleLoginAttempt(Store), HandleTestGetter)

	// DELETE
	users.DELETE("deleteUser", HandleDeleteUser)
}

// @Summary Create a new user
// @tags users
// @Router /api/users/create [post]
func HandleCreateUser(c *gin.Context) {

	// Binds Model and handles validation.
	var json params.CreateUserParams

	if err := c.ShouldBindJSON(&json); err != nil {
		// Handle errors
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect details supplied, please try again."})
		return
	}

	if !helpers.ValidateEmail(json.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
		c.Abort()
		return
	}

	// Validate the password being sent.
	if len(json.Password) <= 7 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password was to short, must be longer than 8 characters."})
		return
	}

	// Validate the password contains at least one letter and capital
	if !helpers.ContainsCapitalLetter(json.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password does not contain a capital letter."})
		return
	}

	// Make sure the password contains at least one special character.
	if !helpers.ContainsSpecialCharacter(json.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The password must contain at least one special character."})
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	// Attempt to create a user.
	insertedId, err := createUser(json.Email, json.Password, json.Type, db.(*gorm.DB))

	if err != nil {
		// Handle the error and or return the context and include a server error status code.
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "The user has been successfully created.",
		"userId":  insertedId,
	})
}

// @Summary Attempt to login using user details
// @tags users
// @Router /api/users/login [post]
func HandleLogin(c *gin.Context) {

	bindJson, _ := c.Get("bindedJson")

	json := bindJson.(params.LoginParams)

	if !helpers.ValidateEmail(json.Email) {
		fmt.Println("A email address was not used to attempt login.")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
		c.Abort()
		return
	}

	// Validate the password being sent.
	if len(json.Password) <= 7 {
		fmt.Println("Password is shorter then 8 characters")
		c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password was to short, must be longer than 8 characters."})
		return
	}

	// Validate the password contains at least one letter and capital
	if !helpers.ContainsCapitalLetter(json.Password) {
		fmt.Println("No Capital letter used.")
		c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password does not contain a capital letter."})
		return
	}

	// Make sure the password contains at least one special character.
	if !helpers.ContainsSpecialCharacter(json.Password) {
		fmt.Println("No special character found.")
		c.JSON(http.StatusBadRequest, gin.H{"message": "The password must contain at least one special character."})
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	session, exists := c.Get("session")

	if !exists {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		return
	}

	userId, outcome, err := loginUser(json.Email, json.Password, db.(*gorm.DB))

	if err != nil {

		// Save changes to our session if an error occurred and we need to abort early..
		if err := Store.Save(c.Request, c.Writer, session.(*sessions.Session)); err != nil {
			fmt.Print(err)
		}

		// Were sending 422 as there is a validation concern.
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		return
	}

	// @todo make this into a map so userid's can be multiple.
	session.(*sessions.Session).Values["userId"] = userId

	if err := Store.Save(c.Request, c.Writer, session.(*sessions.Session)); err != nil {
		fmt.Print(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"attempt": outcome,
		"message": "You have successfully logged into your account.",
	})

}

// @Summary Updates a users details
// @tags users
// @Router /api/users/updateUserDetails [post]
func HandleUpdateUserDetails(c *gin.Context) {
	var json params.UpdateUserParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields, please try again."})
		log.Println(err)
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	outcome, err := updateUser(json.Id, json.Email, json.AccountType, json.FirstName, json.LastName, json.PhoneNumber, json.RecoveryEmail, db.(*gorm.DB))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again."})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": outcome,
	})

}

// @Summary Deletes a user using a user id
// @tags users
// @Router /api/users/deleteUser [delete]
func HandleDeleteUser(c *gin.Context) {
	var json params.DeleteUserParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields, please try again."})
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	outcome, err := deleteUser(json.Id, db.(*gorm.DB))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again."})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": outcome,
	})

}

// @Summary Attempts to get a existing user by id
// @tags users
// @Router /api/users/getUserById [get]
func HandleGetUserById(c *gin.Context) {
	// Were using delete params as it shares the same interface.
	var json params.DeleteUserParams

	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No user ID found, please try again."})
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	outcome, err := getUser(json.Id, db.(*gorm.DB))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully found user",
		"user":    outcome,
	})

}

// @Summary Attempts to get the currently logged in user using there session id.
// @tags users
// @Router /api/users/getCurrentUser [get]
func HandleGetCurrentUser(c *gin.Context) {

	// Get the currently logged int user id.
	userId := c.MustGet("userId")

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	outcome, err := getUser(userId.(uint), db.(*gorm.DB))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully found user",
		"user":    outcome,
	})

}

func HandleTestGetter(c *gin.Context) {

	session, err := c.Get("session")

	if !err {
		fmt.Println("session not found")
	}

	fmt.Printf("%#v\n", session.(*sessions.Session).Values["client"])

	fmt.Printf("%#v\n", session.(*sessions.Session).Values["client"].(ClientProfile).LoginAttempts["test"]["test@liam.pro"])

	// Save changes to our session.
	if err := Store.Save(c.Request, c.Writer, session.(*sessions.Session)); err != nil {
		fmt.Print(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test Ran successfully",
	})

}

func HandleTestPoster(c *gin.Context) {

	sessionValues, err := Store.Get(c.Request, "connect.s.id")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%#v\n", sessionValues)

	c.JSON(http.StatusOK, gin.H{
		"message": "Test Ran successfully",
	})

}
