package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Checks if a user is logged in with a session to the master dashboard;
func ifMasterAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionValues, err := Store.Get(c.Request, "connect.s.id")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to view this."})
			c.Abort()
		}

		// Requires the user to be authorised.
		// @todo make this check for tenancy id also.
		if sessionValues.Values["Authorised"] != true {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to view this."})
			c.Abort()
		}

		// Pass the user id into the handler.
		c.Set("userId", sessionValues.Values["userId"])
	}
}

// Checks if a user is logged in with a session to a tenancy.
func ifAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionValues, err := Store.Get(c.Request, "connect.s.id")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to view this."})
			c.Abort()
		}

		// Requires the user to be authorised.
		// @todo make this check for tenancy id also.
		if sessionValues.Values["Authorised"] != true {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to view this."})
			c.Abort()
		}

		// Pass the user id into the handler.
		c.Set("userId", sessionValues.Values["userId"])
	}
}

// Specific check to see if the current user is also an administrator using there userID
//func checkIfAdmin() gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//	}
//}
