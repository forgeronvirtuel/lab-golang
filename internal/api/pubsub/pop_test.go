package pubsub

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPopHandler(t *testing.T) {
	queue := NewQueue[[]byte]()
	handler := NewPopHandler(queue)

	// Enqueue a sample message
	queue.Enqueue([]byte("sample message"))

	// Create a test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/pop", nil)

	// Call the handler
	handler.HandlePop(c)

	// Check the response
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "sample message" {
		t.Errorf("expected body 'sample message', got '%s'", w.Body.String())
	}

	// Test popping from an empty queue
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/pop", nil)

	// Call the handler again
	handler.HandlePop(c)

	// Check the response for empty queue
	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d for empty queue, got %d", http.StatusNotFound, w.Code)
	}
}
