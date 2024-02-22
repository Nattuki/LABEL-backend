package handler

import (
	"LABEL-backend/user"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Me struct {
	myName       string
	myIconBase64 string
}

func HandleGetMe(c echo.Context) error {
	sess, err := session.Get("LABEL_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}
	accessToken := sess.Values["access_token"]
	log.Println("start" + accessToken.(string))
	return c.JSON(http.StatusOK, Me{
		myName:       user.GetName(accessToken.(string)),
		myIconBase64: user.GetIcon(accessToken.(string)),
	})
}
