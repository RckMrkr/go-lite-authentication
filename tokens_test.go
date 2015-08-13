package main

import (
	// "fmt"
	"bytes"

	"net/http"
	"net/http/httptest"

	"database/sql"
	_ "github.com/erikstmartin/go-testdb"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTokenAndVerifySuccessful(t *testing.T) {
	assert := assert.New(t)

	db, _ := sql.Open("testdb", "")
	router := createRouter(db, []byte("public_key"), []byte("private_key"))

	req, _ := http.NewRequest("POST", "/token/", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if !assert.Equal(200, w.Code) {
		return
	}

	token := w.Body.Bytes()
	req, _ = http.NewRequest("POST", "/token/verify", bytes.NewBuffer(token))
	w = httptest.NewRecorder()

	assert.Equal(200, w.Code)
}

func TestCreateTokenShouldGiveError(t *testing.T) {

}

func TestVerifyTokenShouldReturn400(t *testing.T) {
	// assert := assert.New(t)

	// db, _ := sql.Open("testdb", "")
	// router := createRouter(db, "public_key", "private_key")

	// req, _ := http.NewRequest("POST", "/token/verify", nil)

	// w := httptest.NewRecorder()
	// router.ServeHTTP(w, req)

	// assert.Equal(400, w.Code)
}

func TestVerifyTokenShouldReturn200(t *testing.T) {

}
