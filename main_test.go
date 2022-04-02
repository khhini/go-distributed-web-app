package main

import (
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
	mockUserResp := `{"ping":"pong Hello"}`
	r := SetupRouter()
	r.GET("/:name", IndexHandler)
	req, _ := http.NewRequest("GET", "/Hello", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, mockUserResp, w.Body.String())

}
