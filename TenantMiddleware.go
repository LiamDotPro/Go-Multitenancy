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
		if err := c.Bind(&json); err == nil && len(json.TenancyIdentifier) > 0 {

			var tenantInfo TenantConnectionInformation

			if err := Connection.Where(&TenantConnectionInformation{TenantSubDomainIdentifier: json.TenancyIdentifier}).First(tenantInfo).Error; err != nil {
				fmt.Println("Tenant Identifier passed was not found in database")
				c.AbortWithStatus(400)
			}

			conn, connErr := tenantInfo.getConnection()

			if connErr != nil {
				fmt.Println("Tenant connection could not be made for the request - attempt using json tenancyIdentifier")
			}

			// Set connection into the context for routing
			c.Set("connection", conn)

			c.Next()
		} else {
			// Try and make a connection using the host subdomain
			getTenantConnectionByHost(c.Request.Host, c)
		}
	}
}

func getTenantConnectionByHost(hostStr string, c *gin.Context) {

	connectionInformation, err := getSubdomainInformation(hostStr)

	if err != nil {
		c.AbortWithStatus(400)
	}

	// Make a check for a tenancy identifier being passed by the host as a subdomain identifier

	conn, connErr := connectionInformation.getConnection()

	if connErr != nil {
		fmt.Print("Tenant connection could not be made for the request - attempt using host tenancyIdentifier")
	}

	// Set connection into the context for routing
	c.Set("connection", conn)

	c.Next()
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

	if err := Connection.Where(&TenantConnectionInformation{TenantSubDomainIdentifier: output[0]}).First(tenantInfo).Error; err != nil {
		return TenantConnectionInformation{}, err
	}

	return tenantInfo, nil
}
