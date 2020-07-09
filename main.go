package main

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var VALID_AUTHENTICATIONS = []string{"user", "admin"}

func main() {

	initDb()
	gin.SetMode(gin.ReleaseMode)

	router = gin.Default()

	initializeRoutes()

	router.Run(":8181")

}
