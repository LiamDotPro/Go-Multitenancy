package middleware

import (
	"errors"
	"fmt"
	"github.com/LiamDotPro/Go-Multitenancy/database"
	"github.com/LiamDotPro/Go-Multitenancy/master"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type tenantIdentifierParams struct {
	TenancyIdentifier string `form:"tenant" json:"tenant"`
}

func FindTenancy() gin.HandlerFunc {
	return func(c *gin.Context) {

		var json tenantIdentifierParams

		// Try and find an incoming tenancy identifier on the request
		// @todo check if the conditions actually work for empty string.
		if err := c.Bind(&json); err == nil && len(json.TenancyIdentifier) > 0 {

			val, found := database.TenantMap[json.TenancyIdentifier]

			if !found {
				fmt.Print("Tenant Identifier passed was not found in tenant map")
			}

			conn, connErr := val.GetConnection()

			if connErr != nil {
				fmt.Print("Tenant connection could not be made for the request - attempt using json tenancyIdentifer")
			}

			// Set connection into the context for routing
			c.Set("connection", conn)

			c.Next()
		} else if connectionInformation, err := getSubdomainInformation(c.Request.Host); err == nil {
			// Make a check for a tenancy identifier being passed by the host as a subdomain identifier

			conn, connErr := connectionInformation.GetConnection()

			if connErr != nil {
				fmt.Print("Tenant connection could not be made for the request - attempt using host tenancyIdentifer")
			}

			// Set connection into the context for routing
			c.Set("connection", conn)

			c.Next()
		} else {
			c.AbortWithStatus(400)
		}
	}
}

func getSubdomainInformation(hostStr string) (master.TenantConnectionInformation, error) {

	output := strings.Split(hostStr, ".")

	if len(output) < 2 {
		return master.TenantConnectionInformation{}, errors.New("there was no subdomain present in the string or not enough to split")
	}

	if len(output[0]) <= 0 {
		return master.TenantConnectionInformation{}, errors.New("subdomain was empty")
	}

	val, found := database.TenantMap[output[0]]

	if !found {
		return master.TenantConnectionInformation{}, errors.New("subdomain was not found in tenant collection")
	}

	return val, nil
}

// Checks if a user is logged in with a session.
func IfAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionValues, err := database.Store.Get(c.Request, "connect.s.id")

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