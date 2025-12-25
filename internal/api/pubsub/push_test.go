package pubsub

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandlePush(t *testing.T) {
	gin.SetMode(gin.TestMode)

	queue := NewQueue()
	handler := NewPushHandler(queue)

	r := gin.New()
	r.POST("/push", handler.HandlePush)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/push", bytes.NewBufferString("test message"))

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	if queue.Size() != 1 {
		t.Fatalf("expected queue size 1, got %d", queue.Size())
	}
	value, ok := queue.Dequeue()
	if !ok {
		t.Fatalf("expected successful dequeue")
	}
	if string(value) != "test message" {
		t.Fatalf("expected 'test message', got '%s'", value)
	}
}
