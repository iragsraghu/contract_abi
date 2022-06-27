package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", indexPage)
	router.POST("/api/v1/contract", contractSourceCode)
	router.Run() // listen and serve on
}
