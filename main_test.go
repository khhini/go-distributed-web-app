package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khhini/go-distributed-web-app/handlers"
	"github.com/khhini/go-distributed-web-app/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var objectID primitive.ObjectID
var jwtToken handlers.JWTOutput

func TestIndexHandler(t *testing.T) {
	mockUserResp := `{"ping":"ping"}`
	r := SetupRouter()
	req, err := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockUserResp, w.Body.String())
}

func TestHealthzHandler(t *testing.T) {
	mockUserResp := `{"ping":"ping"}`
	r := SetupRouter()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockUserResp, w.Body.String())
}
func TestSignInJWTHandler(t *testing.T) {
	r := SetupRouter()

	user := models.User{
		Password: "wrongpassword",
		Username: "admin",
	}
	jsonValue, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/signin", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	user = models.User{
		Password: "passadmin",
		Username: "admin",
	}
	jsonValue, _ = json.Marshal(user)
	req, err = http.NewRequest("POST", "/signin", bytes.NewBuffer(jsonValue))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	json.Unmarshal([]byte(w.Body.String()), &jwtToken)

	assert.NotNil(t, jwtToken)
}

func TestRefreshJWTHandler(t *testing.T) {
	r := SetupRouter()

	req, err := http.NewRequest("POST", "/refresh", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `{"error":"token contains an invalid number of segments"}`, w.Body.String())

	req, err = http.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Authorization", jwtToken.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error":"Token is not expired yet"}`, w.Body.String())

	req, err = http.NewRequest("POST", "/refresh", nil)
	req.Header.Set("Authorization", "randominvalidtoken")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `{"error":"token contains an invalid number of segments"}`, w.Body.String())

}

func TestNewRecipeHandler(t *testing.T) {
	r := SetupRouter()

	recipe := models.Recipe{
		Name: "New York Pizza",
	}
	jsonValue, _ := json.Marshal(recipe)
	req, err := http.NewRequest("POST", "/recipes", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, err = http.NewRequest("POST", "/recipes", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "randominvalidtoken")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	req, err = http.NewRequest("POST", "/recipes", bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", jwtToken.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var payload map[string]string
	json.Unmarshal([]byte(w.Body.String()), &payload)
	objectID, _ = primitive.ObjectIDFromHex(payload["recipeID"])

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, objectID)

}

func TestListRecipesHandler(t *testing.T) {
	r := SetupRouter()

	req, err := http.NewRequest("GET", "/recipes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var recipes []models.Recipe
	json.Unmarshal([]byte(w.Body.String()), &recipes)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 493, len(recipes))
}

func TestUpdateRecipeHandler(t *testing.T) {
	r := SetupRouter()

	recipe := models.Recipe{
		ID:   objectID,
		Name: "Gnocchi",
		Ingredients: []string{
			"5 large Idaho potatoes\r",
			"2 eggs\r",
			"3/4 cup grated Parmesan\r",
			"3 1/2 cup all-purpose flour\r",
		},
	}

	jsonValue, _ := json.Marshal(recipe)
	reqFound, err := http.NewRequest("PUT", fmt.Sprintf("/recipes/%s", recipe.ID.Hex()), bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	reqFound, err = http.NewRequest("PUT", fmt.Sprintf("/recipes/%s", recipe.ID.Hex()), bytes.NewBuffer(jsonValue))
	reqFound.Header.Set("Authorization", "randominvalidtoken")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	jsonValue, _ = json.Marshal(recipe)
	reqFound, err = http.NewRequest("PUT", fmt.Sprintf("/recipes/%s", recipe.ID.Hex()), bytes.NewBuffer(jsonValue))
	reqFound.Header.Set("Authorization", jwtToken.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	reqNotFound, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewBuffer(jsonValue))
	reqNotFound.Header.Set("Authorization", jwtToken.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)

}

func TestDeleteRecipeHandler(t *testing.T) {
	r := SetupRouter()

	reqFound, err := http.NewRequest("DELETE", fmt.Sprintf("/recipes/%s", objectID.Hex()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	reqFound, err = http.NewRequest("DELETE", fmt.Sprintf("/recipes/%s", objectID.Hex()), nil)
	reqFound.Header.Set("Authorization", "randominvalidtoken")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	reqFound, err = http.NewRequest("DELETE", fmt.Sprintf("/recipes/%s", objectID.Hex()), nil)
	reqFound.Header.Set("Authorization", jwtToken.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqFound)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	reqNotFound, err := http.NewRequest("DELETE", fmt.Sprintf("/recipes/%s", objectID.Hex()), nil)
	reqNotFound.Header.Set("Authorization", jwtToken.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, reqNotFound)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSearchRecipesHandler(t *testing.T) {
	r := SetupRouter()

	tag := "italian"
	req, err := http.NewRequest("GET", "/recipes/search?tag="+tag, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var recipes []models.Recipe
	json.Unmarshal([]byte(w.Body.String()), &recipes)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	for _, x := range recipes {
		assert.Contains(t, []string(x.Tags), tag)
	}
}
