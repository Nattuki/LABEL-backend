package main

import (
	"LABEL-backend/handler"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var (
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)

func init() {
	handler.ClientID = os.Getenv("TRAQ_CLIENT_ID")
}

func main() {
	// Echoの新しいインスタンスを作成
	e := echo.New()
	e.Use(session.Middleware(store))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello!\n")
	})
	e.GET("/login", handler.HandleLogin)
	e.GET("/callback", handler.HandleCallback)

	err := e.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
