package main

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	configs.ConnectDB()

	routes.UserRoute(r)

	r.Run()
}
