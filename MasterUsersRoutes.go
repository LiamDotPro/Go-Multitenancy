package main

import (
	"fmt"
	"github.com/LiamDotPro/Go-Multitenancy/helpers"
	"github.com/LiamDotPro/Go-Multitenancy/params"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Init
func setupMasterUsersRoutes(router *gin.Engine) {

	users := router.Group("/master/api/users")

	// POST
	users.POST("create", HandleMasterCreateUser)
	users.POST("login", HandleMasterLogin)
	users.POST("updateUserDetails", HandleMasterUpdateUserDetails)
	users.POST("createNewTenant", HandleCreateNewTenant)

	// GET
	users.GET("getUserById", HandleMasterGetUserById)
	users.GET("getCurrentUser", HandleMasterGetCurrentUser)

	// DELETE
	users.DELETE("deleteUser", HandleMasterDeleteUser)
}

// @Summary Create a new user
// @tags master/users
// @Router /master/api/users/create [post]
func HandleMasterCreateUser(c *gin.Context) {

	// Binds Model and handles validation.
	var json params.CreateUserParams

	if err := c.ShouldBindJSON(&json); err != nil {
		// Handle errors
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect details supplied, please try again."})
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

	// Attempt to create a user.
	insertedId, err := createMasterUser(json.Email, json.Password, json.Type)

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
// @tags master/users
// @Router /master/api/users/login [post]
func HandleMasterLogin(c *gin.Context) {

	var json params.LoginParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
		return
	}

	userId, outcome, err := loginMasterUser(json.Email, json.Password)

	if err != nil {
		// Were sending 422 as there is a validation concern.
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		return
	}

	// Setup new session only for host application.
	session, err := Store.New(c.Request, "connect.s.id")

	session.Values["host"] = true
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
// @tags master/users
// @Router /master/api/users/updateUserDetails [post]
func HandleMasterUpdateUserDetails(c *gin.Context) {
	var json params.UpdateUserParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields, please try again."})
		log.Println(err)
		return
	}

	outcome, err := updateMasterUser(json.Id, json.Email, json.AccountType, json.FirstName, json.LastName, json.PhoneNumber, json.RecoveryEmail)

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
// @tags master/users
// @Router /master/api/users/deleteUser [delete]
func HandleMasterDeleteUser(c *gin.Context) {
	var json params.DeleteUserParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields, please try again."})
		return
	}

	outcome, err := deleteMasterUser(json.Id)

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
// @tags master/users
// @Router /master/api/users/getUserById [get]
func HandleMasterGetUserById(c *gin.Context) {
	// Were using delete params as it shares the same interface.
	var json params.DeleteUserParams

	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No user ID found, please try again."})
		return
	}

	outcome, err := getMasterUser(json.Id)

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
// @tags master/users
// @Router /master/api/users/getCurrentUser [get]
func HandleMasterGetCurrentUser(c *gin.Context) {

	// Get the currently logged int user id.
	userId := c.MustGet("userId")

	outcome, err := getMasterUser(userId.(uint))

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

// @Summary Attempts to create a new tenant as a privileged user.
// @tags master/users
// @Router /master/api/users/createNewTenant [Post]
func HandleCreateNewTenant(c *gin.Context) {

	var json params.CreateNewTenantParams

	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No subdomain identifier was found."})
		return
	}

	outcome, err := createNewTenant(json.SubDomainIdentifier)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": outcome,
	})

}
