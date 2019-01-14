package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

type tenantIdentifierParams struct {
	TenancyIdentifier string `form:"tenant" json:"tenant"`
}

func findTenancy() gin.HandlerFunc {
	return func(c *gin.Context) {

		var json tenantIdentifierParams

		// Try and find an incoming tenancy identifier on the request
		// @todo check if the conditions actually work for empty string.
		if err := c.Bind(&json); err == nil && len(json.TenancyIdentifier) > 0 {

			var tenantInfo TenantConnectionInformation

			if err := Connection.Where(&TenantConnectionInformation{TenantSubDomainIdentifier: json.TenancyIdentifier}).First(tenantInfo).Error; err != nil {
				fmt.Print("Tenant Identifier passed was not found in database")
			}

			conn, connErr := tenantInfo.getConnection()

			if connErr != nil {
				fmt.Print("Tenant connection could not be made for the request - attempt using json tenancyIdentifer")
			}

			// Set connection into the context for routing
			c.Set("connection", conn)

			c.Next()
		} else if connectionInformation, err := getSubdomainInformation(c.Request.Host); err == nil {
			// Make a check for a tenancy identifier being passed by the host as a subdomain identifier

			conn, connErr := connectionInformation.getConnection()

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

func getSubdomainInformation(hostStr string) (TenantConnectionInformation, error) {

	output := strings.Split(hostStr, ".")

	if len(output) < 2 {
		return TenantConnectionInformation{}, errors.New("there was no subdomain present in the string or not enough to split")
	}

	if len(output[0]) <= 0 {
		return TenantConnectionInformation{}, errors.New("subdomain was empty")
	}

	var tenantInfo TenantConnectionInformation

	if err := Connection.Where(&TenantConnectionInformation{TenantSubDomainIdentifier: hostStr}).First(tenantInfo).Error; err != nil {
		return TenantConnectionInformation{}, err
	}

	return tenantInfo, nil
}
