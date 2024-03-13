package handler

import (
	"LABEL-backend/user"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

type Message struct {
	MessageId   string    `json:"messageId" db:"message_id"`
	CreatorName string    `json:"creatorName" db:"creator_name"`
	Title       string    `json:"title" db:"title"`
	Comment     string    `json:"comment" db:"comment"`
	Url         string    `json:"url" db:"url"`
	CreatedOn   time.Time `json:"createdOn" db:"created_on"`
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
	message.MessageId = "msg" + xid.New().String()
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

	return c.NoContent(http.StatusOK)
}

func (h *dbHandler) HandleGetMessage(c echo.Context) error {
	messageId := c.Param("messageid")

	var message Message

	err := h.db.Get(&message, "SELECT * From messages WHERE message_id = ?", messageId)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get the message from the database")
	}

	return c.JSON(http.StatusOK, message)
}

func (h *dbHandler) HandleGetMessages(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusNotFound, "invalid path parameter")
	}

	name := c.QueryParam("name")

	var message Message
	var messages []Message
	var messagesToSend []Message
	var rows *sqlx.Rows

	if name == "" {
		rows, err = h.db.Queryx("SELECT * From messages ORDER BY created_on DESC")
	} else {
		rows, err = h.db.Queryx("SELECT * From messages WHERE creator_name = ? ORDER BY created_on DESC", name)
	}
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get messages from the database")
	}

	for rows.Next() {
		err = rows.StructScan(&message)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "failed to scan the next row")
		}
		messages = append(messages, message)
	}

	messagesSize := len(messages)
	start := 10 * (page - 1)
	if start > (messagesSize-1) || start < 0 {
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

func (h *dbHandler) HandleCountPages(c echo.Context) error {
	name := c.QueryParam("name")
	res := struct {
		Count int `json:"count"`
	}{0}

	var err error
	if name == "" {
		err = h.db.Get(&res.Count, "SELECT COUNT(message_id) FROM messages")
	} else {
		h.db.Get(&res.Count, "SELECT COUNT(*) FROM messages WHERE creator_name = ?", name)
	}
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to connect with the database")
	}

	res.Count = (res.Count-1)/10 + 1
	return c.JSON(http.StatusOK, res)
}

func (h *dbHandler) HandleDeleteMessage(c echo.Context) error {
	messageId := c.Param("id")

	sess, err := session.Get("LABEL_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}
	accessToken := sess.Values["access_token"]
	userName := user.GetName(accessToken.(string))

	var creatorName string
	err = h.db.Get(&creatorName, "SELECT creator_name FROM messages WHERE message_id = ?", messageId)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to connect with the database")
	}

	if creatorName != userName {
		return c.String(http.StatusUnauthorized, "you are not authorized to delete the message")
	}

	_, err = h.db.Exec("DELETE FROM labels WHERE message_id = ?", messageId)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to delete the message")
	}

	_, err = h.db.Exec("DELETE FROM messages WHERE message_id = ?", messageId)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to delete the message")
	}

	return c.NoContent(http.StatusOK)
}
