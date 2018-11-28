package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func findTenancy() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println(&c)

		// Set example variable
		c.Set("connection", "12345")

		c.Next()
	}
}
