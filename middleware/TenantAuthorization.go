package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wader/gormstore"
	"net/http"
)

// Checks if a user is logged in with a session to a tenancy.
func IfAuthorized(Store *gormstore.Store) gin.HandlerFunc {
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
