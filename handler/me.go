package handler

import (
	"LABEL-backend/user"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type Me struct {
	MyId         string `json:"MyId"`
	MyName       string `json:"MyName"`
	MyIconBase64 string `json:"MyIconBase64"`
	IsVisitor    bool   `json:"IsVisitor"`
}

func HandleGetMe(c echo.Context) error {
	sess, err := session.Get("LABEL_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}

	accessToken := sess.Values["access_token"]

	if user.IsValidToken(accessToken.(string)) {
		return c.JSON(http.StatusOK, &Me{
			MyId:         user.GetId(accessToken.(string)),
			MyName:       user.GetName(accessToken.(string)),
			MyIconBase64: user.GetIcon(accessToken.(string)),
			IsVisitor:    false,
		})
	} else {
		return c.JSON(http.StatusOK, &Me{
			MyId:         "",
			MyName:       "",
			MyIconBase64: "",
			IsVisitor:    true,
		})
	}
}
