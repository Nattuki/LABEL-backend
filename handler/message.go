package handler

import (
	"LABEL-backend/user"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

type Message struct {
	MessageId   string    `json:"-" db:"message_id"`
	CreatorName string    `json:"-" db:"creator_name"`
	Title       string    `json:"title" db:"title"`
	Comment     string    `json:"comment" db:"comment"`
	Url         string    `json:"url" db:"url"`
	CreatedOn   time.Time `json:"-" db:"created_on"`
}

func (h *dbHandler) HandleSendMessage(c echo.Context) error {
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
	message.CreatedOn = time.Now()

	_, err = h.db.Exec("INSERT INTO messages (message_id, creator_name, title, comment, url, created_on) VALUES (?, ?, ?, ?, ?, ?)",
		message.MessageId,
		message.CreatorName,
		message.Title,
		message.Comment,
		message.Url,
		message.CreatedOn,
	)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to insert into the database")
	}

	log.Println(*message)
	return c.String(http.StatusOK, "OK!")
}

func (h *dbHandler) HandleGetMessage(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusNotFound, "invalid path parameter")
	}

	var message Message
	var messages []Message
	var messagesToSend []Message

	rows, err := h.db.Queryx("SELECT * From messages ORDERBY created_on")
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get messages from the database")
	}

	for rows.Next() {
		err := rows.StructScan(&message)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "failed to scan the next row")
		}
		messages = append(messages, message)
	}

	messagesSize := len(messages)
	start := 10 * (page - 1)
	if start > (messagesSize - 1) {
		return c.String(http.StatusNotFound, "invalid path parameter")
	}
	var end int
	if 10*page <= messagesSize {
		end = 10 * page
	} else {
		end = messagesSize
	}

	messagesToSend = messages[start:end]

	return c.JSON(http.StatusOK, messagesToSend)
}
