package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var playerCollection = configs.GetCollection(configs.DB, "players")
var validate = validator.New()

func CreatePlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		player := models.Player{}
		defer cancel()

		if err := c.ShouldBindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		if validationErr := validate.Struct(&player); validationErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": validationErr.Error(),
			})
			return
		}

		newPlayer := models.Player{
			Id:       primitive.NewObjectID(),
			Name:     player.Name,
			Region:   player.Region,
			Position: player.Position,
		}

		_, err := playerCollection.InsertOne(ctx, newPlayer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, newPlayer)
	}
}

func GetPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		playerId := c.Param("playerId")
		player := models.Player{}
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(playerId)

		err := playerCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&player)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, player)
	}
}

func EditPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		playerId := c.Param("playerId")
		player := models.Player{}
		defer cancel()

		ObjId, _ := primitive.ObjectIDFromHex(playerId)

		if err := c.ShouldBindJSON(&player); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": err.Error(),
			})
			return
		}

		validationErr := validate.Struct(&player)
		if validationErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": validationErr.Error(),
			})
			return
		}

		update := bson.M{
			"name":     player.Name,
			"region":   player.Region,
			"position": player.Position,
		}

		_player := models.Player{}
		err := playerCollection.FindOne(ctx, bson.M{"_id": ObjId}).Decode(&_player)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Player with specified ID not found!",
			})
			return
		}

		result, err := playerCollection.UpdateOne(ctx, bson.M{"_id": ObjId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}

		updatedPlayer := models.Player{}
		if result.MatchedCount == 1 {
			err := playerCollection.FindOne(ctx, bson.M{"_id": ObjId}).Decode(&updatedPlayer)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": err.Error(),
				})
				return
			}
		}
		c.JSON(http.StatusOK, updatedPlayer)
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Player with specified ID not found!",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Player successfully deleted!"})
	}
}

func GetAllPlayers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		players := []models.Player{}
		defer cancel()

		result, err := playerCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}
		defer result.Close(ctx)
		for result.Next(ctx) {
			siglePlayer := models.Player{}
			if err := result.Decode(&siglePlayer); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": err.Error(),
				})
			}
			players = append(players, siglePlayer)
		}
		c.JSON(http.StatusOK, players)
	}
}
