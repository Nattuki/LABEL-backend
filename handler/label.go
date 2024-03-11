package handler

import (
	"LABEL-backend/user"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

type Label struct {
	LabelId     string    `json:"labelId" db:"label_id"`
	MessageId   string    `json:"messageId" db:"message_id"`
	Content     string    `json:"content" db:"content"`
	CreatorName string    `json:"creatorName" db:"creator_name"`
	CreatedOn   time.Time `json:"createdOn" db:"createdOn"`
}

func (h *dbHandler) HandleSendLabel(c echo.Context) error {
	label := new(Label)
	err := c.Bind(label)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get the label")
	}

	sess, err := session.Get("LABEL_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}
	accessToken := sess.Values["access_token"]

	label.CreatorName = user.GetName(accessToken.(string))
	label.LabelId = "label" + xid.New().String()
	label.CreatedOn = time.Now()

	_, err = h.db.Exec("INSERT INTO labels (label_id, message_id, content, creator_name, created_on) VALUES (?, ?, ?, ?, ?)",
		label.LabelId,
		label.MessageId,
		label.Content,
		label.CreatorName,
		label.CreatedOn,
	)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to insert into the database")
	}

	log.Println(*label)
	return c.NoContent(http.StatusOK)
}
