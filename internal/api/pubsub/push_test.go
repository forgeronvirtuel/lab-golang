package pubsub

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandlePush(t *testing.T) {
	queue := NewQueue()
	handler := NewPushHandler()

	// Create a test context with a sample payload
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	payload := []byte("test message")
	c.Request, _ = http.NewRequest("POST", "/push", bytes.NewBuffer(payload))

	// Call the handler
	handler.HandlePush(c)

	// Check the response
	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Verify that the message was enqueued
	if queue.Size() != 1 {
		t.Errorf("expected queue size 1, got %d", queue.Size())
	}
	value, ok := queue.Dequeue()
	if !ok {
		t.Errorf("expected successful dequeue")
	}
	if string(value) != "test message" {
		t.Errorf("expected 'test message', got '%s'", value)
	}
}
