package middleware

import (
	"errors"
	"fmt"
	"github.com/LiamDotPro/Go-Multitenancy/tenants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strings"
)

type tenantIdentifierParams struct {
	TenancyIdentifier string `form:"tenant" json:"tenant"`
}

func FindTenancy(Connection *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var json tenantIdentifierParams

		// Try and find an incoming tenancy identifier on the request
		if err := c.Bind(&json); err == nil && len(json.TenancyIdentifier) > 0 {

			var tenantInfo tenants.TenantConnectionInformation

			if err := Connection.Where(&tenants.TenantConnectionInformation{TenantSubDomainIdentifier: json.TenancyIdentifier}).First(&tenantInfo).Error; err != nil {
				fmt.Println("Tenant Identifier passed was not found in database")
				c.AbortWithStatus(400)
			}

			conn, connErr := tenantInfo.GetConnection()

			if connErr != nil {
				fmt.Println("Tenant connection could not be made for the request - attempt using json tenancyIdentifier")
			}

			// Set connection into the context for routing
			c.Set("connection", conn)

			// Set tenancy Identifier into params
			c.Set("tenantIdentifier", json.TenancyIdentifier)

			c.Next()
		} else {
			// Try and make a connection using the host subdomain
			getTenantConnectionByHost(c.Request.Host, c, Connection)
		}
	}
}

func getTenantConnectionByHost(hostStr string, c *gin.Context, Connection *gorm.DB) {

	connectionInformation, tenantString, err := getSubdomainInformation(hostStr, Connection)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	// Make a check for a tenancy identifier being passed by the host as a subdomain identifier
	conn, connErr := connectionInformation.GetConnection()

	if connErr != nil {
		fmt.Println("Tenant connection could not be made for the request - attempt using host tenancyIdentifier")
	}

	// Set connection into the context for routing
	c.Set("connection", conn)

	// Set tenancy Identifier into params
	c.Set("tenantIdentifier", tenantString)

	c.Next()
}

func getSubdomainInformation(hostStr string, Connection *gorm.DB) (TenantConnectionInfo tenants.TenantConnectionInformation, tenantIdentifier string, err error) {

	output := strings.Split(hostStr, ".")

	if len(output) < 2 {
		return tenants.TenantConnectionInformation{}, "", errors.New("there was no subdomain present in the string or not enough to split")
	}

	if len(output[0]) <= 0 {
		return tenants.TenantConnectionInformation{}, "", errors.New("subdomain was empty")
	}

	var tenantInfo tenants.TenantConnectionInformation

	if err := Connection.Where(&tenants.TenantConnectionInformation{TenantSubDomainIdentifier: output[0]}).First(&tenantInfo).Error; err != nil {
		return tenants.TenantConnectionInformation{}, "", errors.New("tenancy identifier not found in database")
	}

	return tenantInfo, output[0], nil
}
