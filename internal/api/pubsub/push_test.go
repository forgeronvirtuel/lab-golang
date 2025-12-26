package pubsub

import (
	"bytes"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// buildFrameData creates a properly formatted frame with channel and data
func buildFrameData(channel string, data []byte) *bytes.Buffer {
	buf := new(bytes.Buffer)
	chLen := uint8(len(channel))
	binary.Write(buf, binary.BigEndian, chLen)
	buf.WriteString(channel)
	dataLen := uint32(len(data))
	binary.Write(buf, binary.BigEndian, dataLen)
	buf.Write(data)
	return buf
}

func TestHandlePush(t *testing.T) {
	gin.SetMode(gin.TestMode)

	queue := NewQueue()
	handler := NewPushHandler(queue)

	r := gin.New()
	r.POST("/push", handler.HandlePush)

	tests := []struct {
		name           string
		channel        string
		message        []byte
		expectedStatus int
		expectEnqueued bool
	}{
		{
			name:           "valid message",
			channel:        "test",
			message:        []byte("test message"),
			expectedStatus: http.StatusCreated,
			expectEnqueued: true,
		},
		{
			name:           "empty message",
			channel:        "events",
			message:        []byte{},
			expectedStatus: http.StatusCreated,
			expectEnqueued: true,
		},
		{
			name:           "large message",
			channel:        "data",
			message:        make([]byte, 1024*1024), // 1 MiB
			expectedStatus: http.StatusCreated,
			expectEnqueued: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue = NewQueue() // Reset queue for each test
			handler.Queue = queue

			w := httptest.NewRecorder()
			body := buildFrameData(tt.channel, tt.message)
			req := httptest.NewRequest("POST", "/push", body)

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectEnqueued {
				if queue.Size() != 1 {
					t.Errorf("expected queue size 1, got %d", queue.Size())
				}
				value, ok := queue.Dequeue()
				if !ok {
					t.Error("expected successful dequeue")
				}
				if !bytes.Equal(value, tt.message) {
					t.Errorf("expected %q, got %q", tt.message, value)
				}
			} else {
				if queue.Size() != 0 {
					t.Errorf("expected queue size 0, got %d", queue.Size())
				}
			}
		})
	}
}

func TestHandlePush_InvalidFrames(t *testing.T) {
	gin.SetMode(gin.TestMode)

	queue := NewQueue()
	handler := NewPushHandler(queue)

	r := gin.New()
	r.POST("/push", handler.HandlePush)

	tests := []struct {
		name           string
		body           []byte
		expectedStatus int
	}{
		{
			name:           "empty body",
			body:           []byte{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "truncated channel length",
			body:           []byte{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "zero channel length",
			body:           []byte{0x00}, // chLen = 0
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "channel too large",
			body: func() []byte {
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.BigEndian, uint8(255))
				buf.Write(make([]byte, 255))
				binary.Write(buf, binary.BigEndian, uint32(100))
				return buf.Bytes()
			}(),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue = NewQueue() // Reset queue
			handler.Queue = queue

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/push", bytes.NewReader(tt.body))

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d, body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if queue.Size() != 0 {
				t.Errorf("expected queue size 0, got %d", queue.Size())
			}
		})
	}
}

func TestHandlePush_MaxBodySize(t *testing.T) {
	gin.SetMode(gin.TestMode)

	queue := NewQueue()
	handler := NewPushHandler(queue)

	r := gin.New()
	r.POST("/push", handler.HandlePush)

	// Create a frame that exceeds 50 MiB
	channel := "test"
	dataLen := uint32(51 << 20) // 51 MiB

	buf := new(bytes.Buffer)
	chLen := uint8(len(channel))
	binary.Write(buf, binary.BigEndian, chLen)
	buf.WriteString(channel)
	binary.Write(buf, binary.BigEndian, dataLen)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/push", buf)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	if queue.Size() != 0 {
		t.Errorf("expected queue size 0, got %d", queue.Size())
	}
}
