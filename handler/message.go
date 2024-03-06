package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Message struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func HandleMessage(c echo.Context) error {
	message := new(Message)
	err := c.Bind(message)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get the message")
	}
	log.Println(*message)
	return c.String(http.StatusOK, "OK!")
}
