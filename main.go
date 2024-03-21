package main

import (
	handler "LABEL-backend/handler"
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/srinathgs/mysqlstore"
)

type Template struct {
	templates *template.Template
}

var (
	tpl = `<!DOCTYPE html">
	<html>
	<head prefix="og:http://ogp.me/ns#">
	<title>転送</title>
	<meta http-equiv="refresh" content="3;url=https://jurassic-design.com/web/redirect-sample-code/">
	<meta property="og:url" content="https://label.trap.show/api" />
	<meta property="og:type" content="article" />
	<meta property="og:title" content="リダイレクトページだよ" />
	<meta property="og:description" content="テストテストテスト" />
	<meta property="og:site_name" content="LABEL" />
	<meta property="og:image" content="https://jurassic-design.com/wp-content/uploads/redirect-sample-code-tittle.jpg" />
	</head>
	
	<body>
	
	<p>〇〇〇〇のホームページは移転しました。</p>
	<p>新しいホームページは、<a href="https://jurassic-design.com/web/redirect-sample-code/">https://jurassic-design.com/web/redirect-sample-code/</a>になります。
	</p>
	
	<p>
	お手数をおかけしますが、「お気に入り」の変更をお願いします。
	なお、このページは、3秒後に、新しいホームページに自動転送されます。
	</p>
	
	</body>
	
	</html>
	`
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

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

	h := handler.NewHandler(db)

	t := &Template{
		templates: template.Must(template.New("hello").Parse(tpl)),
	}
	e.Renderer = t

	e.GET("/", h.HandleRenderMessage)
	e.GET("/me", handler.HandleGetMe)
	e.GET("/loginpath", handler.HandleGetOAuthUrl)
	e.GET("/gettoken", handler.HandleGetToken)
	e.GET("/message/:id", h.HandleGetMessage)
	e.GET("/message/get/:page", h.HandleGetMessages)
	e.GET("/message/countPages", h.HandleCountPages)
	e.GET("/label/get/:messageid", h.HandleGetLabel)
	e.POST("/message/send", h.HandleSendMessage)
	e.POST("/label/send", h.HandleSendLabel)
	e.DELETE("/message/:id", h.HandleDeleteMessage)
	e.DELETE("/label/:id", h.HandleDeleteLabel)
	e.DELETE("/logout", handler.HandleLogout)
	err = e.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
