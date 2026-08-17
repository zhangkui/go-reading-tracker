package httpapi

import (
	"github.com/zhangkui/go-reading-tracker/internal/reading"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAndListBooks(t *testing.T) {
	server := NewServer(reading.NewService())
	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"Go","author":"A","isbn":"10","total_pages":100}`))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest(http.MethodGet, "/books?page=1&page_size=10", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"total":1`) {
		t.Fatalf("list status = %d body=%s", response.Code, response.Body.String())
	}
}
