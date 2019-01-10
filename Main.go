package main

import (
	"github.com/LiamDotPro/Go-Multitenancy/database"
	"github.com/LiamDotPro/Go-Multitenancy/tenants/users"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"io"
	"os"
)

func main() {
	// Show Swagger pages
	// @todo This needs to be made a dev only process by env var
	//open("http://localhost:8000/swagger/index.html")

	// Configure port
	port := ":" + os.Getenv("PORT")

	if port == ":" {
		port = ":8000"
	}

	// Start database services and load master database.
	database.StartDatabaseServices()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// init router
	router := gin.Default()

	// Setting up our routes on the router.

	// Users
	users.SetupUsersRoutes(router)

	// Add routing for swag @todo make this development only using envs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Starting the router instance
	if err := router.Run(port); err != nil {

	}
}
