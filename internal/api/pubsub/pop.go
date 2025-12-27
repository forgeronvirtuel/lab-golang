package pubsub

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PopHandler struct {
	Queue *Queue[[]byte]
}

func NewPopHandler(queue *Queue[[]byte]) *PopHandler {
	return &PopHandler{Queue: queue}
}

// HandlePop processes pop requests and dequeues messages.
func (h *PopHandler) HandlePop(c *gin.Context) {
	// Dequeue a message
	value, ok := h.Queue.Dequeue()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "no messages in queue"})
		return
	}

	// Respond with the message
	c.Data(http.StatusOK, "application/octet-stream", value)
}
