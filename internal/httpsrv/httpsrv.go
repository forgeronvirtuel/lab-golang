package httpsrv

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

func StartHTTPServer(srv *http.Server, wg *sync.WaitGroup) {
	// Start server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}

func StopHTTPServer(srv *http.Server, timeout time.Duration) error {
	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
