package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func UserAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "Failed to get the session.")
		}

		if sess.Values["userName"] == nil {
			c.Set("userName", "visitor")
		} else {
			c.Set("userName", sess.Values["userName"].(string))
		}
		log.Println(sess.Values)

		return next(c)
	}
}
