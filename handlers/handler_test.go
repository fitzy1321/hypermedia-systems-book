package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := Index()
	handler.ServeHTTP(w, req)

	res := w.Result()
	// body := w.Body.String()

	assert.Equal(t, http.StatusSeeOther, res.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", res.Header.Get("Content-Type"))
	// assert.Contains(t, body, "Hello, TestUser")
}
