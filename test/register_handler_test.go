package test

import (
	"basic_server/server/model"
	"basic_server/test/service"
	"basic_server/test/service/database"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterAttemptWithEmptyRequestPayload(t *testing.T) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", nil)
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRegisterAttemptWithEmptyEmailField(t *testing.T) {
	requestPayload, _ := json.Marshal(map[string]string{
		"email":    "",
		"password": "test password",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestFailedRegisterAttemptIfUserAlreadyExists(t *testing.T) {
	cleaner := database.Cleaner(service.TestServer().DatabaseDriver())

	defer cleaner.CleanUp()

	service.TestServer().DatabaseDriver().Create(
		&model.User{
			Email:    "test1@test.com",
			Password: "test",
		},
	)

	requestPayload, _ := json.Marshal(map[string]string{
		"email":    "test1@test.com",
		"password": "test password",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSuccessfulRegisterAttempt(t *testing.T) {
	cleaner := database.Cleaner(service.TestServer().DatabaseDriver())

	defer cleaner.CleanUp()

	requestPayload, _ := json.Marshal(map[string]string{
		"email":     "test2@test.com",
		"password":  "test password",
		"full_name": "test test",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSuccessfulRegisterAttemptWithoutFullNameField(t *testing.T) {
	cleaner := database.Cleaner(service.TestServer().DatabaseDriver())

	defer cleaner.CleanUp()

	requestPayload, _ := json.Marshal(map[string]string{
		"email":    "test3@test.com",
		"password": "test password",
	})

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestPayload))
	req.Header.Add("Content-Type", "application/json")

	service.TestServer().Engine().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}