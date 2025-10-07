package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexHandler(t *testing.T) {
	// Arrange http request, Recorder, and handler
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := Index()

	// Act
	handler.ServeHTTP(w, req)
	res := w.Result()
	// body := w.Body.String()

	// Assert
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", res.Header.Get("Content-Type"))
	// assert.Contains(t, body, "Hello, TestUser")
}

// func TestGetContactHandler(t *testing.T) {
// 	// Arrange http request, recorder, and handler
// 	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
// 	w := httptest.NewRecorder()

// 	db, err := dbmodels.GetAndSetupSqlite3DB("")
// 	assert.NoError(t, err, "Error raised trying to make in memory sqlite db")
// 	contact_template := template.New("Hello, Contact Page!")

// 	handler := Contacts(db, contact_template)

// 	// Act
// 	handler.ServeHTTP(w, req)
// 	res := w.Result()

// 	// Assert
// 	assert.Equal(t, http.StatusOK, res.StatusCode)
// }
