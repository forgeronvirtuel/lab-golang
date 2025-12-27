package pubsub

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PushHandler struct {
	Queue *Queue[*Frame]
}

func NewPushHandler(q *Queue[*Frame]) *PushHandler {
	return &PushHandler{Queue: q}
}

const maxBody = int64(50 << 20) // 50 MiB
const maxChannelLen = uint8(255)

// HandlePush processes incoming push requests and enqueues messages.
func (h *PushHandler) HandlePush(c *gin.Context) {

	// Read the header
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBody)
	frame, br, err := ReadFrameHeader(c.Request.Body, maxChannelLen, 50<<20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Read the message data
	data := make([]byte, frame.DataLen)
	if _, err := br.Read(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read message data"})
		return
	}

	frame.Data = data

	// Enqueue the message
	h.Queue.Enqueue(&frame)

	// Respond with success
	c.Status(http.StatusCreated)
}
