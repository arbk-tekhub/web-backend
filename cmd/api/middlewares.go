package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) requireBasicAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, plaintextPassword, ok := c.Request.BasicAuth()
		if !ok {
			app.basicAuthenticationRequired(c)
			c.Abort()
			return
		}

		if app.config.basicAuth.username != username {
			app.basicAuthenticationRequired(c)
			c.Abort()
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(app.config.basicAuth.hashedPassword), []byte(plaintextPassword))
		if err != nil {
			switch {
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
				app.basicAuthenticationRequired(c)
				c.Abort()
				return
			default:
				app.internalServerErrorResponse(c, err)
				c.Abort()
				return
			}

		}
		c.Next()
	}
}
