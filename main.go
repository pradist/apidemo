package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/pradist/apidemo/auth"
	"github.com/pradist/apidemo/todo"
)

func main() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("tokenz", auth.AccessToken("==signature=="))
	protect := r.Group("", auth.Protect([]byte("==signature==")))

	handler := todo.NewTodoHandler(db)
	protect.POST("/todos", handler.NewTask)

	r.Run()
}
