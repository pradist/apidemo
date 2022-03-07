package todo

import (
	"log"
	"net/http"
	"time"
)

type Todo struct {
	Title     string `json:"text"`
	ID        uint   `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Todo) TableName() string {
	return "todos"
}

type storer interface {
	New(*Todo) error
}

type TodoHandler struct {
	store storer
}

func NewTodoHandler(store storer) *TodoHandler {
	return &TodoHandler{store: store}
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	TransactionID() string
	Audience() string
}

func (t *TodoHandler) NewTask(c Context) {
	var todo Todo
	// if err := c.ShouldBindJSON(&todo); err != nil {
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if todo.Title == "sleep" {
		// transactionID := c.Request.Header.Get("TransactionID")
		transactionID := c.TransactionID()
		// aud, _ := c.Get("aud")
		aud := c.Audience()
		log.Println(transactionID, aud, "not allow")
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "not allow",
		})
		return
	}

	if err := t.store.New(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"ID": todo.ID,
	})
}
