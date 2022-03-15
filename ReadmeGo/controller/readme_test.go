package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// router.POST("/readme", createReadme)
// router.GET("/readme/:id", getReadme)
// router.PUT("/readme/:id/header", addHeader)
// router.PUT("/readme/:id/code", addCode)
// router.PUT("/readme/:id/blockquote", addBlockquote)
// router.PUT("/readme/:id/link", addLink)
// router.PUT("/readme/:id/image", addImage)
// router.POST("/readme/:id/file", createReadmeFile)

func TestGetReadme(t *testing.T) {
	router := setupRouter()
	r := httptest.NewRecorder()
	w := httptest.NewRecorder()

	req1, _ := http.NewRequest("POST", "/readme?name=1", nil)
	router.ServeHTTP(r, req1)

	emptyReadme := string(`[""]`)

	req, _ := http.NewRequest("GET", "/readme/1", nil)
	router.ServeHTTP(w, req)

	require.JSONEq(t, emptyReadme, w.Body.String())
}

func TestGetReadmeReturnsNotFound(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/readme/1", nil)
	router.ServeHTTP(w, req)

	require.JSONEq(t, string(`{"message": "could not find readme"}`), w.Body.String())
}

func TestAddHeader(t *testing.T) {
	// router := setupRouter()
	// w := httptest.NewRecorder()
}
