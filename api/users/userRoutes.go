package users

import (
	_ "../docs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Init
func setupUsersRoutes(router *gin.Engine) {

	users := router.Group("/api/users")

	// POST
	users.POST("create", HandleCreateUser)
	users.POST("login", HandleLogin)
	users.POST("updateUserDetails", ifAuthorized(), HandleUpdateUserDetails)
	users.POST("assignUserCompany", HandleUpdateUserCompany)

	// GET
	users.GET("getUserById", ifAuthorized(), HandleGetUserById)
	users.GET("getAllSystemUsers", HandleGetAllSystemUsers)
	users.GET("getAllCompanyUsers", HandleGetAllUsersByCompanyId)
	users.GET("getCurrentUser", ifAuthorized(), HandleGetCurrentUser)

	// DELETE
	users.DELETE("deleteUser", ifAuthorized(), checkIfAdmin(), HandleDeleteUser)
}

/**
Model binding types.
 */

type CreateUserParams struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Type     int    `form:"type" json:"type"`
}

type UpdateUserParams struct {
	Id            uint   `form:"id" json:"id" binding:"required"`
	Email         string `form:"email" json:"email"`
	AccountType   int    `form:"accountType" json:"accountType"`
	FirstName     string `form:"firstName" json:"firstName"`
	LastName      string `form:"lastName" json:"lastName"`
	PhoneNumber   string `form:"phoneNumber" json:"phoneNumber"`
	RecoveryEmail string `form:"recoveryEmail" json:"recoveryEmail"`
	Contractor    bool   `form:"contractor" json:"contractor"`
}

type UserCompanyUpdateParams struct {
	Id        uint `form:"id" json:"id" binding:"required"`
	CompanyId int  `form:"companyId" json:"companyId" binding:"required"`
}

type LoginParams struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type DeleteUserParams struct {
	Id uint `form:"id" json:"id" binding:"required"`
}

type UserCompanySearchParams struct {
	CompanyId uint `form:"companyId" json:"companyId" binding:"required"`
}

// @Summary Create a new user
// @tags users
// @Router /api/users/create [post]
func HandleCreateUser(c *gin.Context) {

	// Binds Model and handles validation.
	var json CreateUserParams

	if err := c.ShouldBindJSON(&json); err != nil {
		// Handle errors
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect details supplied, please try again."})
		return
	}

	// Attempt to create a user.
	insertedId, err := createUser(json.Email, json.Password, json.Type)

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

	var json LoginParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
		return
	}

	userId, outcome, err := loginUser(json.Email, json.Password)

	if err != nil {
		// Were sending 422 as there is a validation concern.
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		return
	}

	// Setup new session.
	session, err := Store.New(c.Request, "connect.s.id")

	session.Values["Authorised"] = true
	session.Values["userId"] = userId

	Store.Save(c.Request, c.Writer, session)

	c.JSON(http.StatusOK, gin.H{
		"attempt": outcome,
		"message": "You have successfully logged into your account.",
	})

}

// @Summary Updates a users details
// @tags users
// @Router /api/users/updateUserDetails [post]
func HandleUpdateUserDetails(c *gin.Context) {
	var json UpdateUserParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields, please try again."})
		log.Println(err)
		return
	}

	outcome, err := updateUser(json.Id, json.Email, json.AccountType, json.FirstName, json.LastName, json.PhoneNumber, json.RecoveryEmail, json.Contractor)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again."})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": outcome,
	})

}

// @Summary Updates a users company assignment
// @tags users
// @Router /api/users/assignUserCompany [post]
func HandleUpdateUserCompany(c *gin.Context) {
	var json UserCompanyUpdateParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields, please try again."})
		log.Println(err)
		return
	}

	outcome, err := updateUserCompany(json.Id, json.CompanyId)

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
	var json DeleteUserParams

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields, please try again."})
		return
	}

	outcome, err := deleteUser(json.Id)

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
	var json DeleteUserParams

	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No user ID found, please try again."})
		return
	}

	outcome, err := getUser(json.Id)

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

	outcome, err := getUser(userId.(uint))

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

// @Summary Attempts to get all available users independent of company
// @tags users
// @Router /api/users/getAllSystemUsers [get]
func HandleGetAllSystemUsers(c *gin.Context) {

	// @todo add super admin check

	outcome, err := getUserList()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved all users.",
		"user":    outcome,
	})

}

// @Summary Attempts to get all users associated with a company.
// @tags users
// @Router /api/users/getAllCompanyUsers [get]
func HandleGetAllUsersByCompanyId(c *gin.Context) {
	var json UserCompanySearchParams

	if err := c.Bind(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No company ID found, please try again."})
		return
	}

	outcome, err := getUserListByCompanyId(json.CompanyId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong while trying to process that, please try again.", "error": err.Error()})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully retrieved all company users.",
		"user":    outcome,
	})

}
