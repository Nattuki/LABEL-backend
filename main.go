package main

import (
	handler "LABEL-backend/handler"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/srinathgs/mysqlstore"
)

func init() {
	handler.ClientID = os.Getenv("TRAQ_CLIENT_ID")
}

func main() {

	// データーベースの設定
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	conf := mysql.Config{
		User:                 os.Getenv("NS_MARIADB_USER"),
		Passwd:               os.Getenv("NS_MARIADB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("NS_MARIADB_HOSTNAME") + ":" + os.Getenv("NS_MARIADB_PORT"),
		DBName:               os.Getenv("NS_MARIADB_DATABASE"),
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	// データベースに接続
	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	// usersテーブルが存在しなかったら、usersテーブルを作成する
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (Username VARCHAR(255) PRIMARY KEY, HashedPass VARCHAR(255))")
	if err != nil {
		log.Fatal(err)
	}

	// セッションの情報を記憶するための場所をデータベース上に設定
	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		log.Fatal(err)
	}

	// Echoの新しいインスタンスを作成
	e := echo.New()
	e.Use(session.Middleware(store))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello!\n")
	})
	h := handler.NewHandler(db)

	e.GET("/me", handler.HandleGetMe)
	e.GET("/loginpath", handler.HandleGetOAuthUrl)
	e.GET("/gettoken", handler.HandleGetToken)
	e.GET("/message/:messageid", h.HandleGetMessage)
	e.GET("/message/get/:page", h.HandleGetMessages)
	e.GET("/message/countPages", h.HandleCountPages)
	e.GET("/label/get/:messageid", h.HandleGetLabel)
	e.POST("/message/send", h.HandleSendMessage)
	e.POST("/label/send", h.HandleSendLabel)
	e.DELETE("/message/delete/:id", h.HandleDeleteMessage)
	e.DELETE("/label/delete/:id", h.HandleDeleteLabel)
	e.DELETE("/logout", handler.HandleLogout)
	err = e.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
