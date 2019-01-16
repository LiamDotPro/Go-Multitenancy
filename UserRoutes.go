package main

import (
	"fmt"
	"github.com/LiamDotPro/Go-Multitenancy/middleware"
	"github.com/LiamDotPro/Go-Multitenancy/params"
	"github.com/gin-gonic/gin"
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
	users.GET("testGetter", HandleMasterLoginAttempt(Store), HandleTestGetter)

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

	var json params.LoginParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
		return
	}

	// Get the database object from the connection.
	db, _ := c.Get("connection")

	userId, outcome, err := loginUser(json.Email, json.Password, db.(*gorm.DB))

	if err != nil {
		// Were sending 422 as there is a validation concern.
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		return
	}

	// Setup new session.
	session, err := Store.New(c.Request, "connect.s.id")

	session.Values["Authorised"] = true
	session.Values["userId"] = userId

	if err := Store.Save(c.Request, c.Writer, session); err != nil {
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

	sessionValues, err := Store.Get(c.Request, "connect.s.id")

	if err != nil {
		fmt.Println(err)
	}

	hostValues := sessionValues.Values["host"].(HostProfile)

	fmt.Printf("%#v\n", hostValues)

	loginAttemptValues := hostValues.LoginAttempts["liam@liams.pro"]

	fmt.Printf("%#v\n", loginAttemptValues.LoginAttempts)
	fmt.Println(loginAttemptValues.LoginAttempts)

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
