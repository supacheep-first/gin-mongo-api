package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var playerCollection *mongo.Collection = configs.GetCollection(configs.DB, "players")
var validate = validator.New()

func CreatePlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var player models.Player
		defer cancel()

		if err := c.BindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, responses.PlayerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&player); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PlayerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newPlayer := models.Player{
			Id:       primitive.NewObjectID(),
			Name:     player.Name,
			Region:   player.Region,
			Position: player.Position,
		}

		result, err := playerCollection.InsertOne(ctx, newPlayer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.PlayerResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		playerId := c.Param("playerId")
		var player models.Player
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(playerId)

		err := playerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&player)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.PlayerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": player}})
	}
}

func EditPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		playerId := c.Param("playerId")
		var player models.Player
		defer cancel()

		ObjId, _ := primitive.ObjectIDFromHex(playerId)

		if err := c.BindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, responses.PlayerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&player); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.PlayerResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{
			"name":     player.Name,
			"region":   player.Region,
			"position": player.Position,
		}

		result, err := playerCollection.UpdateOne(ctx, bson.M{"_id": ObjId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedPlayer models.Player
		if result.MatchedCount == 1 {
			err := playerCollection.FindOne(ctx, bson.M{"_id": ObjId}).Decode(&updatedPlayer)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}
		c.JSON(http.StatusOK, responses.PlayerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedPlayer}})
	}
}

func DeletePlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		playerId := c.Param("playerId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(playerId)

		result, err := playerCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.PlayerResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Player with specified ID not found!"}})
			return
		}
		c.JSON(http.StatusOK, responses.PlayerResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Player successfully deleted!"}})
	}
}

func GetAllPlayers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var players []models.Player
		defer cancel()

		result, err := playerCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		defer result.Close(ctx)
		for result.Next(ctx) {
			var siglePlayer models.Player
			if err := result.Decode(&siglePlayer); err != nil {
				c.JSON(http.StatusInternalServerError, responses.PlayerResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}
			players = append(players, siglePlayer)
		}
		c.JSON(http.StatusOK, responses.PlayerResponse{Status: http.StatusOK, Message: "players", Data: map[string]interface{}{"data": players}})
	}

}
