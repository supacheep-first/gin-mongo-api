package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func PlayerRoute(router *gin.Engine) {
	router.POST("/player", controllers.CreatePlayer())
	router.GET("/player/:playerId", controllers.GetPlayer())
	router.PUT("/player/:playerId", controllers.EditPlayer())
	router.DELETE("/player/:playerId", controllers.DeletePlayer())
	router.GET("/players", controllers.GetAllPlayers())
}
