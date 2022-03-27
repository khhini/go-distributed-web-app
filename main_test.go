package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestIndexHandler(t *testing.T) {
	mockUserResp := `{"ping":"pong"}`
	r := SetupRouter()
	r.GET("/", IndexHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockUserResp, w.Body.String())
}

func TestNewRecipeHandler(t *testing.T) {
	r := SetupRouter()
	r.POST("/recipes", NewRecipeHandler)

	recipe := Recipe{
		Name: "New York Pizza",
	}
	jsonValue, _ := json.Marshal(recipe)
	req, _ := http.NewRequest("POST", "/recipes", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListRecipesHandler(t *testing.T) {
	r := SetupRouter()
	r.GET("/recipes", ListRecipesHandler)

	req, _ := http.NewRequest("GET", "/recipes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var recipes []Recipe
	json.Unmarshal([]byte(w.Body.String()), &recipes)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 493, len(recipes))
}
