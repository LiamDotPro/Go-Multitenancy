package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type tenantIdentifierParams struct {
	TenancyIdentifier string `form:"tenant" json:"tenant" binding:"required"`
}

func findTenancy() gin.HandlerFunc {
	return func(c *gin.Context) {

		var json tenantIdentifierParams

		// Try and find an incoming tenancy identifier on the request
		c.Bind(&json)

		fmt.Printf("%+v\n", json)

		// Set example variable
		c.Set("connection", &User{})

		c.Next()
	}
}
