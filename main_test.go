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
