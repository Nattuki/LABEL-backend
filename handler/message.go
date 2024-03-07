package handler

import (
	"LABEL-backend/user"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

type Message struct {
	MessageId   string `json:"-" db:"message_id"`
	CreatorName string `json:"-" db:"creator_name"`
	Title       string `json:"title" db:"title"`
	Url         string `json:"url" db:"url"`
}

func (h *dbHandler) HandleMessage(c echo.Context) error {
	message := new(Message)
	err := c.Bind(message)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get the message")
	}

	sess, err := session.Get("LABEL_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}
	accessToken := sess.Values["access_token"]
	message.CreatorName = user.GetName(accessToken.(string))
	message.MessageId = xid.New().String()

	log.Println(*message)
	return c.String(http.StatusOK, "OK!")
}
