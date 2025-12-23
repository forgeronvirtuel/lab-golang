package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/forgeronvirtuel/lab-golang/internal/httpsrv"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var pubsubCmd = &cobra.Command{
	Use:   "pubsub",
	Short: "Start an HTTP pub/sub server",
	Long:  `Start an HTTP server that provides pub/sub functionality.`,
	Example: `  # Start the server on default port 8080
  lab-golang pubsub

  # Start the server on a specific port
  lab-golang pubsub --port 3000

  # Start the server on a specific host and port
  lab-golang pubsub --host 0.0.0.0 --port 8080
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Configure Gin
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()

		// Define routes
		router.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		router.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		// Configure HTTP server
		addr := fmt.Sprintf("%s:%s", serverHost, serverPort)
		srv := &http.Server{
			Addr:    addr,
			Handler: router,
		}

		log.Printf("Starting server on %s\n", addr)
		var wg sync.WaitGroup
		httpsrv.StartHTTPServer(srv, &wg)

		// Wait for interrupt signal for graceful shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("Shutting down server...")
		if err := httpsrv.StopHTTPServer(srv, 5*time.Second); err != nil {
			log.Fatalf("Server forced to shutdown: %v", err)
		}

		wg.Wait()
		log.Println("Exiting application, bye!")
	},
}

func init() {
	rootCmd.AddCommand(pubsubCmd)

	pubsubCmd.Flags().StringVarP(&serverPort, "port", "p", "8080", "Port to listen on")
	pubsubCmd.Flags().StringVarP(&serverHost, "host", "H", "localhost", "Host to bind to")
}
