package pubsub

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PushHandler struct {
	Queue *Queue
}

func NewPushHandler(q *Queue) *PushHandler {
	return &PushHandler{Queue: q}
}

// HandlePush processes incoming push requests and enqueues messages.
func (h *PushHandler) HandlePush(c *gin.Context) {
	// Read the request body
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}

	// Enqueue the message
	h.Queue.Enqueue(body)

	// Respond with success
	c.Status(http.StatusCreated)
}
