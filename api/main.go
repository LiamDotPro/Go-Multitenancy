package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"io"
	"os"
	"os/exec"
	"runtime"
	"github.com/liamdotpro/database"
	"net/http"
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

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// init router
	router := gin.Default()

	// Setting up our routes on the router.
	// Users
	setupUsersRoutes(router)

	// Add routing for swag @todo make this development only using envs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Starting the router instance
	router.Run(port)
}

// Checks if a user is logged in with a session.
func ifAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {

		sessionValues, err := database.Store.Get(c.Request, "connect.s.id")

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to view this."})
			c.Abort()
		}

		// Requires the user to be authorised.
		if sessionValues.Values["Authorised"] != true {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to view this."})
			c.Abort()
		}

		// Pass the user id into the handler.
		c.Set("userId", sessionValues.Values["userId"])
	}
}

// Specific check to see if the current user is also an administrator using there userID
func checkIfAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// Attempts to migrate tables by using file specific functions
func migrateTables() {
	database.Connection.AutoMigrate(&User{})
}

// Helper function that allows us to open a browser dependant on your OS
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
