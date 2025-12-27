package pubsub

import (
	"net/http"

	"github.com/forgeronvirtuel/lab-golang/internal/httpsrv"
	"github.com/gin-gonic/gin"
)

type PushHandler struct {
	Queue *Queue[[]byte]
}

func NewPushHandler(q *Queue[[]byte]) *PushHandler {
	return &PushHandler{Queue: q}
}

const maxBody = int64(50 << 20) // 50 MiB
const maxChannelLen = uint8(255)

// HandlePush processes incoming push requests and enqueues messages.
func (h *PushHandler) HandlePush(c *gin.Context) {

	// Read the header
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBody)
	header, br, err := httpsrv.ReadFrameHeader(c.Request.Body, maxChannelLen, 50<<20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Read the message data
	data := make([]byte, header.DataLen)
	if _, err := br.Read(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read message data"})
		return
	}

	// Enqueue the message
	h.Queue.Enqueue(data)

	// Respond with success
	c.Status(http.StatusCreated)
}
