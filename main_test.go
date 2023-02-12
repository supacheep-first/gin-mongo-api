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

	body := []byte(`{
		"name": "Test Lionel Messi",
		"region": "Argentina",
		"position": "FW"
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/player", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NotNil(t, w.Body)
	resp := gin.H{}
	err := json.Unmarshal([]byte(w.Body.String()), &resp)

	id = resp["id"].(string)
	assert.Nil(t, err)
	assert.Equal(t, models.Player{
		Name:     "Test Lionel Messi",
		Region:   "Argentina",
		Position: "FW",
	}, models.Player{
		Name:     resp["name"].(string),
		Region:   resp["region"].(string),
		Position: resp["position"].(string),
	})
}

func TestGetPlayer(t *testing.T) {
	r := gin.Default()
	routes.PlayerRoute(r)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/player/%s", id), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)

	resp := gin.H{}
	err := json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Nil(t, err)
	reqId, _ := primitive.ObjectIDFromHex(id)
	respId, _ := primitive.ObjectIDFromHex(resp["id"].(string))
	assert.Equal(t, models.Player{
		Id:       reqId,
		Name:     "Test Lionel Messi",
		Region:   "Argentina",
		Position: "FW",
	}, models.Player{
		Id:       respId,
		Name:     resp["name"].(string),
		Region:   resp["region"].(string),
		Position: resp["position"].(string),
	})
}

func TestUpdatePlayer(t *testing.T) {
	r := gin.Default()
	routes.PlayerRoute(r)

	body := []byte(`{
        "name": "Test Lionel Messi Updated",
        "region": "Argentina",
		"position": "FW"
	}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/player/%s", id), bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, w.Body)

	resp := gin.H{}
	err := json.Unmarshal([]byte(w.Body.String()), &resp)

	assert.Nil(t, err)
	reqId, _ := primitive.ObjectIDFromHex(id)
	respId, _ := primitive.ObjectIDFromHex(resp["id"].(string))
	assert.Equal(t, models.Player{
		Id:       reqId,
		Name:     "Test Lionel Messi Updated",
		Region:   "Argentina",
		Position: "FW",
	}, models.Player{
		Id:       respId,
		Name:     resp["name"].(string),
		Region:   resp["region"].(string),
		Position: resp["position"].(string),
	})
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
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	assert.Nil(t, err)
	assert.Equal(t, gin.H{
		"message": "Player successfully deleted!",
	}, resp)
}
