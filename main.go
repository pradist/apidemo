package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/pradist/apidemo/auth"
	"github.com/pradist/apidemo/todo"
)

func main() {

	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variables %s", err)
	}

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

	r.GET("tokenz", auth.AccessToken(os.Getenv("SIGN")))
	protect := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	handler := todo.NewTodoHandler(db)
	protect.POST("/todos", handler.NewTask)

	r.Run()
}
