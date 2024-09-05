package main

import (
	interpreter "backend/interpreter"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ExecuteRequest struct {
	Code string `json:"code" binding:"required"`
}

type ExecuteResponse struct {
	Received string `json:"received"`
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
	router.POST("/run-code", func(c *gin.Context) {
		var req ExecuteRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ExecuteResponse{
				Received: err.Error(),
			})
			return
		}
		code := req.Code
		fmt.Println("Received code:", code)
		tokens, err := interpreter.Lex(code)
		if err != nil {
			c.JSON(http.StatusBadRequest, ExecuteResponse{
				Received: err.Error(),
			})
			return
		}
		stack, err := interpreter.Parse(tokens)
		if err != nil {
			c.JSON(http.StatusBadRequest, ExecuteResponse{
				Received: err.Error(),
			})
			return
		}
		output, err := interpreter.Execute(stack)
		if err != nil {
			c.JSON(http.StatusBadRequest, ExecuteResponse{
				Received: err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, ExecuteResponse{
			Received: output,
		})
	})
	fmt.Println("Server running http://localhost:8080")
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
