package cmd

var (
	outputFile string
	numRows    int
	filePath   string
	separator  string
	showFirst  int
	hasHeader  bool
	groupByCol int
	validate   bool
	filters    []string
	// Pub/Sub server flags
	serverPort string
	serverHost string
)
