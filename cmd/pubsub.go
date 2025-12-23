package cmd

import (
	"fmt"

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
		fmt.Println("Not implemented yet")
	},
}

func init() {
	rootCmd.AddCommand(pubsubCmd)

	pubsubCmd.Flags().StringVarP(&serverPort, "port", "p", "8080", "Port to listen on")
	pubsubCmd.Flags().StringVarP(&serverHost, "host", "H", "localhost", "Host to bind to")
}
