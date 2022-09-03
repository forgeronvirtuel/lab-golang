package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Jsonifiable interface {
	ToJson() ([]byte, error)
}

type HttpErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Details map[string]string `json:"details"`
}

type RouteResponse struct {
	Path string `json:"path"`
}

func (r *RouteResponse) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

func writeInternalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	errResponse := HttpErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: fmt.Sprintf("Error happened: %s", err),
		Details: map[string]string{
			"err": err.Error(),
		},
	}
	data, _ := json.Marshal(errResponse)
	w.Write(data)
}

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create data response and convert it to JSON
	rr := RouteResponse{
		Path: r.URL.Path,
	}
	data, err := rr.ToJson()
	if err != nil {
		writeInternalError(w, err)
		return
	}

	// write response
	w.Write(data)
}

func main() {
	addr := os.Args[1]

	// Configure the http server
	s := &http.Server{
		Addr:           addr,
		Handler:        &Handler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		log.Printf("Shutting down server due to interrupt signal")
		if err := s.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	// Starting the server and waiting for a signal stop
	log.Printf("Listening at %s", addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-idleConnsClosed

	log.Printf("Server closed")
}
