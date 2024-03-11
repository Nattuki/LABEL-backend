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
	JumpTime    int       `json:"jumpTime" db:"jump_time"`
	CreatorName string    `json:"creatorName" db:"creator_name"`
	CreatedOn   time.Time `json:"createdOn" db:"created_on"`
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

	_, err = h.db.Exec("INSERT INTO labels (label_id, message_id, content, jump_time, creator_name, created_on) VALUES (?, ?, ?, ?, ?, ?)",
		label.LabelId,
		label.MessageId,
		label.Content,
		label.JumpTime,
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

func (h *dbHandler) HandleGetLabel(c echo.Context) error {
	messageId := c.Param("messageid")

	var label Label
	var labels []Label

	rows, err := h.db.Queryx("SELECT * From labels WHERE message_id = ? ORDER BY created_on DESC", messageId)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get labels from the database")
	}

	for rows.Next() {
		err = rows.StructScan(&label)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "failed to scan the next row")
		}
		labels = append(labels, label)
	}

	return c.JSON(http.StatusOK, labels)
}
