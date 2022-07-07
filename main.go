package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// create a gin router with default middleware:
	router := gin.Default()

	// add a route to the gin router for the root path
	router.GET("/", indexPage)

	// add a route to the gin router for the contract source code
	router.POST("/api/v1/contract", contractSourceCode)

	// start the gin router
	router.Run()
}
