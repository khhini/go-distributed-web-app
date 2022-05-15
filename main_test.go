package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestIndexHandler(t *testing.T) {
	mockUserResp := `{"ping":"ping"}`
	r := SetupRouter()
	r.GET("/", IndexHandler)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockUserResp, w.Body.String())
}

func TestHealthzHandler(t *testing.T) {
	mockUserResp := `{"ping":"ping"}`
	r := SetupRouter()
	r.GET("/healthz", IndexHandler)
	req, _ := http.NewRequest("GET", "/healthz", nil)
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

func TestUpdateRecipeHandler(t *testing.T) {
	r := SetupRouter()
	r.PUT("/recipes/:id", UpdateRecipeHandler)

	recipe := Recipe{
		ID:   "c0283p3d0cvuglq85lpg",
		Name: "Gnocchi",
		Ingredients: []string{
			"5 large Idaho potatoes\r",
			"2 eggs\r",
			"3/4 cup grated Parmesan\r",
			"3 1/2 cup all-purpose flour\r",
		},
	}

	jsonValue, _ := json.Marshal(recipe)
	reqFound, _ := http.NewRequest("PUT", "/recipes/"+recipe.ID, bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)

	assert.Equal(t, http.StatusOK, w.Code)

	reqNotFound, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestDeleteRecipeHandler(t *testing.T) {
	r := SetupRouter()
	r.DELETE("/recipes/:id", DeleteRecipeHandler)
	recipeID := "c0283p3d0cvuglq85lpg"
	reqFound, _ := http.NewRequest("DELETE", "/recipes/"+recipeID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)

	assert.Equal(t, http.StatusOK, w.Code)

	reqNotFound, _ := http.NewRequest("DELETE", "/recipes/"+recipeID, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSearchRecipesHandler(t *testing.T) {
	r := SetupRouter()
	r.GET("/recipes/search", SearchRecipesHandler)
	tag := "italian"
	req, _ := http.NewRequest("GET", "/recipes/search?tag="+tag, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var recipes []Recipe
	json.Unmarshal([]byte(w.Body.String()), &recipes)
	assert.Equal(t, http.StatusOK, w.Code)

	for _, x := range recipes {
		assert.Contains(t, []string(x.Tags), tag)
	}

}
