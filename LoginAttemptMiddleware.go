package main

import (
	"fmt"
	"github.com/LiamDotPro/Go-Multitenancy/params"
	"github.com/gin-gonic/gin"
	"github.com/wader/gormstore"
	"net/http"
	"time"
)

type ClientProfile struct {
	LoginAttempts    map[string]*loginAttempt // Key is email address
	AuthorizationMap map[string]bool          // Key is tenant identifier
}

type HostProfile struct {
	LoginAttempts        map[string]*loginAttempt // Key is used email address
	LastLoginAttemptTime time.Time
	AuthorizedTime       time.Time
	UserId               uint
	Authorized           bool
}

type loginAttempt struct {
	LastLoginAttemptTime time.Time
	LoginAttempts        uint
}

func newHostProfile() HostProfile {
	h := HostProfile{}
	h.LoginAttempts = make(map[string]*loginAttempt)
	h.Authorized = false
	return h
}

// Checks if a user is logged in with a session to the master dashboard;
func HandleMasterLoginAttempt(Store *gormstore.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check our parameters out.
		var json params.LoginParams

		// Abort if we don't have the correct variables to begin with.
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
			c.Abort()
			return
		}

		// Try and get a session.
		sessionValues, err := Store.Get(c.Request, "connect.s.id")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
			c.Abort()
			return
		}

		// Check to see if a new session is found.
		if sessionValues.ID == "" {
			// Setup new session with empty profiles.
			session, err := Store.New(c.Request, "connect.s.id")

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
				c.Abort()
				return
			}

			// Host profile requires little setup
			session.Values["host"] = newHostProfile()

			// Client profile requires no setup
			session.Values["client"] = ClientProfile{}

			if err := Store.Save(c.Request, c.Writer, session); err != nil {
				fmt.Print(err)
			}

		} else {
			// Profile was already found
			h := sessionValues.Values["host"].(HostProfile)

			// Check to see if were already authorized with the host
			if h.Authorized {
				c.JSON(http.StatusOK, gin.H{"message": "Something went wrong..", "status": "Already Authorized with the application."})
				c.Abort()
				return
			}

			// Check if the email used is already in our login attempts.
			loginAttemptsFound, found := h.LoginAttempts[json.Email]

			if !found {
				// email has not been used to login add a new entry
				h.LoginAttempts[json.Email] = &loginAttempt{LoginAttempts: 1, LastLoginAttemptTime: time.Now().UTC()}

				if err := Store.Save(c.Request, c.Writer, sessionValues); err != nil {
					fmt.Print(err)
				}
				return
			}

			// Check to see if login attempts exceeds 3 attempts
			if found && loginAttemptsFound.LoginAttempts > 2 {
				// Check to see if last login attempt was over half an hour ago
				if time.Now().Sub(loginAttemptsFound.LastLoginAttemptTime).Minutes() > 30 {
					// reset login attempts to have 2 more.
					loginAttemptsFound.LoginAttempts = 1
					loginAttemptsFound.LastLoginAttemptTime = time.Now().UTC()
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"message": "You have been locked out for too many attempts to login..", "status": "locked out", "timeLeft": 30 - time.Now().Sub(loginAttemptsFound.LastLoginAttemptTime).Minutes()})
					c.Abort()
					return
				}
			}

			if found && loginAttemptsFound.LoginAttempts <= 2 {
				// increase login attempt count
				loginAttemptsFound.LoginAttempts++
				// replace last attempt date
				loginAttemptsFound.LastLoginAttemptTime = time.Now().UTC()

				if err := Store.Save(c.Request, c.Writer, sessionValues); err != nil {
					fmt.Print(err)
				}
			}

		}

	}

}
