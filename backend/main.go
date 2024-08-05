package main

import (
	interpreter "backend/interpreter"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	router.POST("/run-code", func(c *gin.Context) {
		code := c.PostForm("code")
		tokens, err := interpreter.Lex(code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		tree, err := interpreter.Parse(tokens)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		interpreter.Traverse(tree)
		c.JSON(http.StatusOK, gin.H{
			"received": code,
		})
	})
	fmt.Println("Server running http://localhost:8080")
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
