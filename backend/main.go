package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	fmt.Println("Server running http://localhost:8080")
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
