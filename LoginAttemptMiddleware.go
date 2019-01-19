package main

import (
	"fmt"
	"github.com/LiamDotPro/Go-Multitenancy/helpers"
	"github.com/LiamDotPro/Go-Multitenancy/params"
	"github.com/gin-gonic/gin"
	"github.com/wader/gormstore"
	"net/http"
	"time"
)

type ClientProfile struct {
	LoginAttempts    map[string]map[string]*loginAttempt // Key is email address
	AuthorizationMap map[string]uint                     // Key is tenant identifier
}

type HostProfile struct {
	LoginAttempts        map[string]*loginAttempt // Key is used email address
	LastLoginAttemptTime time.Time
	AuthorizedTime       time.Time
	UserId               uint
	Authorized           uint
}

type loginAttempt struct {
	LastLoginAttemptTime time.Time
	LoginAttempts        uint
}

func newHostProfile() HostProfile {
	h := HostProfile{}
	h.LoginAttempts = make(map[string]*loginAttempt)
	h.Authorized = 0
	return h
}

// Checks if a user is logged in with a session to the master dashboard
func HandleMasterLoginAttempt(Store *gormstore.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Try and get a session.
		sessionValues, err := Store.Get(c.Request, "connect.s.id")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
			c.Abort()
			return
		}

		// Check to see if the user is already authorized..
		if sessionValues.ID != "" {

			p := sessionValues.Values["host"].(HostProfile)

			if p.Authorized == 1 {
				c.JSON(http.StatusOK, gin.H{
					"outcome": "Already Authorized",
					"message": "user already authorized with application.",
				})
				c.Abort()
				return
			}
		}

		// Check our parameters out.
		var json params.LoginParams

		// Abort if we don't have the correct variables to begin with.
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
			fmt.Println("Can't bind request variables for login")
			c.Abort()
			return
		}

		if !helpers.ValidateEmail(json.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
			fmt.Println("Email is not in a valid format.")
			c.Abort()
			return
		}

		// Abort if the passed password is not correct.

		// Validate the password being sent.
		if len(json.Password) <= 7 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password was to short, must be longer than 8 characters."})
			c.Abort()
			return
		}

		// Validate the password contains at least one letter and capital
		if !helpers.ContainsCapitalLetter(json.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password does not contain a capital letter."})
			c.Abort()
			return
		}

		// Make sure the password contains at least one special character.
		if !helpers.ContainsSpecialCharacter(json.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The password must contain at least one special character."})
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
			c.Abort()
			return
		}

		c.Set("bindedJson", json)

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

			return
		} else {
			// Profile was already found
			h := sessionValues.Values["host"].(HostProfile)

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
					return
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
				return
			}

		}

	}

}

// Checks if a user is logged in with a session to the client dashboard
func HandleLoginAttempt(Store *gormstore.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Try and get tenancy identifier
		tenantIdentifier, found := c.Get("tenantIdentifier")

		if !found {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
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

		// Check to see if the user is already authorized..
		if sessionValues.ID != "" {

			p := sessionValues.Values["client"].(ClientProfile)

			authorizationEntry := p.AuthorizationMap[tenantIdentifier.(string)]

			if authorizationEntry == 1 {
				c.JSON(http.StatusOK, gin.H{
					"outcome": "Already Authorized",
					"message": "user already authorized with application.",
				})
				c.Abort()
				return
			}
		}

		// Check our parameters out.
		var json params.LoginParams

		// Abort if we don't have the correct variables to begin with.
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
			fmt.Println("Can't bind request variables for login")
			c.Abort()
			return
		}

		if !helpers.ValidateEmail(json.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email or Password provided are incorrect, please try again."})
			fmt.Println("Email is not in a valid format.")
			c.Abort()
			return
		}

		// Abort if the passed password is not correct.

		// Validate the password being sent.
		if len(json.Password) <= 7 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password was to short, must be longer than 8 characters."})
			c.Abort()
			return
		}

		// Validate the password contains at least one letter and capital
		if !helpers.ContainsCapitalLetter(json.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The specified password does not contain a capital letter."})
			c.Abort()
			return
		}

		// Make sure the password contains at least one special character.
		if !helpers.ContainsSpecialCharacter(json.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "The password must contain at least one special character."})
			c.Abort()
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong.."})
			c.Abort()
			return
		}

		c.Set("bindedJson", json)

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

			return
		} else {
			// Profile was already found
			h := sessionValues.Values["client"].(ClientProfile)

			// Attempt to find tenant entry in login attempts.
			tenantMap, found := h.LoginAttempts[tenantIdentifier.(string)]

			if !found {
				// Create a new entry for the tenant entry in map, also create login attempt
				tenantMap = map[string]*loginAttempt{}
				tenantMap[json.Email] = &loginAttempt{LoginAttempts: 1, LastLoginAttemptTime: time.Now().UTC()}
				return
			}

			// Check if the email used is already in our login attempts.
			loginAttemptsFound, found := tenantMap[json.Email]

			if !found {
				// email has not been used to login add a new entry
				loginAttemptsFound = &loginAttempt{LoginAttempts: 1, LastLoginAttemptTime: time.Now().UTC()}

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
					return
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
				return
			}

		}

	}

}
