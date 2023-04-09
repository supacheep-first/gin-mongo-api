package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-mongo-api/models"
	"gin-mongo-api/routes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var id string

func TestCreatePlayer(t *testing.T) {
	r := gin.Default()
	routes.PlayerRoute(r)

	playerData := models.Player{
		Name:     "Test Lionel Messi",
		Region:   "Argentina",
		Position: "FW",
	}
	body, _ := json.Marshal(playerData) // convert struct to json

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/player", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotNil(t, w.Body)
	resp := models.Player{}
	err := json.Unmarshal(w.Body.Bytes(), &resp) // convert json to struct

	id = resp.Id.Hex()
	assert.Nil(t, err)
	assert.Equal(t, playerData, models.Player{
		Name:     resp.Name,
		Region:   resp.Region,
		Position: resp.Position,
	}) // compare the created player
}

func TestUpdatePlayer(t *testing.T) {
	r := gin.Default()
	routes.PlayerRoute(r)

	playerData := models.Player{
		Name:     "Test Lionel Messi Updated",
		Region:   "Argentina",
		Position: "FW",
	}
	body, _ := json.Marshal(playerData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/player/%s", id), bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)

	resp := models.Player{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	playerData.Id, _ = primitive.ObjectIDFromHex(id) // convert string to ObjectID
	resp.Id, _ = primitive.ObjectIDFromHex(resp.Id.Hex())
	assert.Equal(t, playerData, resp) // compare the updated player
}

func TestDeletePlayer(t *testing.T) {
	r := gin.Default()
	routes.PlayerRoute(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/player/%s", id), nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)

	resp := gin.H{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, gin.H{
		"message": "Player successfully deleted!",
	}, resp)
}
