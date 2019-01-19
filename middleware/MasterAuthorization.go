package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wader/gormstore"
	"net/http"
)

// Checks if a user is logged in with a session to the master dashboard;
func IfMasterAuthorized(Store *gormstore.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionValues, err := Store.Get(c.Request, "connect.s.id")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to view this."})
			c.Abort()
		}

		if sessionValues.Values["Authorised"] != true {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to view this."})
			c.Abort()
		}

		// Pass the user id into the handler.
		c.Set("userId", sessionValues.Values["userId"])
	}
}
