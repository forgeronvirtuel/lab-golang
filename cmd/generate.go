package cmd

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var (
	outputFile string
	numRows    int
)

// Données boursières factices
var (
	symbols = []string{
		"AAPL", "MSFT", "GOOGL", "AMZN", "META", "TSLA", "NVDA", "JPM", "V", "JNJ",
		"WMT", "PG", "MA", "HD", "DIS", "NFLX", "ADBE", "CRM", "ORCL", "INTC",
		"CSCO", "PFE", "KO", "PEP", "NKE", "MCD", "ABT", "TMO", "COST", "AVGO",
		"BA", "IBM", "GE", "CAT", "AMD", "QCOM", "TXN", "SBUX", "INTU", "PYPL",
	}
	exchanges  = []string{"NYSE", "NASDAQ", "EURONEXT", "LSE", "TSE"}
	sectors    = []string{"Technology", "Healthcare", "Financial", "Consumer", "Energy", "Industrial", "Utilities", "Materials", "Real Estate"}
	orderTypes = []string{"Market", "Limit", "Stop", "Stop-Limit"}
	tradeTypes = []string{"Buy", "Sell"}
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a CSV file with stock market data",
	Long: `Generate a large CSV file containing realistic stock market transaction data.
The file includes multiple columns such as trade details, OHLC prices, volume, and market metrics.`,
	Example: `  lab-golang generate --output stock_data.csv --rows 1000000
  lab-golang generate -o data.csv -r 10000000`,
	Run: func(cmd *cobra.Command, args []string) {
		generateStockData(outputFile, numRows)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "stock_data.csv", "Output CSV file path")
	generateCmd.Flags().IntVarP(&numRows, "rows", "r", 1000000, "Number of rows to generate")
}

func generateStockData(output string, rows int) {
	fmt.Printf("Generating %d rows of CSV data...\n", rows)
	start := time.Now()

	file, err := os.Create(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Écrire l'en-tête
	header := []string{
		"TradeID",
		"Timestamp",
		"Symbol",
		"Exchange",
		"Sector",
		"TradeType",
		"OrderType",
		"Quantity",
		"Price",
		"TotalValue",
		"OpenPrice",
		"ClosePrice",
		"HighPrice",
		"LowPrice",
		"Volume",
		"MarketCap",
		"PERatio",
		"DividendYield",
		"Beta",
		"52WeekHigh",
		"52WeekLow",
		"ChangePercent",
	}

	if err := writer.Write(header); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing header: %v\n", err)
		os.Exit(1)
	}

	// Seed pour la génération aléatoire
	rand.Seed(time.Now().UnixNano())

	// Générer les lignes avec timestamps séquentiels
	baseDate := time.Date(2020, 1, 1, 9, 30, 0, 0, time.UTC) // Ouverture marché

	// Stocker les prix de base pour chaque symbole pour simuler des variations réalistes
	stockPrices := make(map[string]float64)
	for _, symbol := range symbols {
		stockPrices[symbol] = float64(rand.Intn(500)+10) + rand.Float64()
	}

	for i := 1; i <= rows; i++ {
		// Timestamp avec progression temporelle
		timestamp := baseDate.Add(time.Duration(i) * time.Second)

		symbol := symbols[rand.Intn(len(symbols))]
		exchange := exchanges[rand.Intn(len(exchanges))]
		sector := sectors[rand.Intn(len(sectors))]
		tradeType := tradeTypes[rand.Intn(len(tradeTypes))]
		orderType := orderTypes[rand.Intn(len(orderTypes))]

		// Prix avec variation réaliste autour du prix de base
		basePrice := stockPrices[symbol]
		variation := (rand.Float64() - 0.5) * basePrice * 0.02 // ±1% variation
		price := basePrice + variation
		stockPrices[symbol] = price // Mise à jour du prix

		quantity := rand.Intn(1000) + 1
		totalValue := price * float64(quantity)

		// Prix OHLC (Open, High, Low, Close)
		openPrice := basePrice * (1 + (rand.Float64()-0.5)*0.01)
		closePrice := price
		highPrice := price * (1 + rand.Float64()*0.015)
		lowPrice := price * (1 - rand.Float64()*0.015)

		// Volume de transactions
		volume := rand.Intn(10000000) + 100000

		// Capitalisation boursière (en millions)
		marketCap := float64(rand.Intn(500000)+1000) * 1000000

		// Ratio P/E
		peRatio := rand.Float64()*50 + 5

		// Rendement du dividende (%)
		dividendYield := rand.Float64() * 5

		// Beta (volatilité)
		beta := rand.Float64()*2 + 0.5

		// Plus haut et plus bas sur 52 semaines
		week52High := basePrice * (1 + rand.Float64()*0.5)
		week52Low := basePrice * (1 - rand.Float64()*0.3)

		// Variation en pourcentage
		changePercent := (price - openPrice) / openPrice * 100

		row := []string{
			strconv.Itoa(i),
			timestamp.Format("2006-01-02 15:04:05"),
			symbol,
			exchange,
			sector,
			tradeType,
			orderType,
			strconv.Itoa(quantity),
			fmt.Sprintf("%.2f", price),
			fmt.Sprintf("%.2f", totalValue),
			fmt.Sprintf("%.2f", openPrice),
			fmt.Sprintf("%.2f", closePrice),
			fmt.Sprintf("%.2f", highPrice),
			fmt.Sprintf("%.2f", lowPrice),
			strconv.Itoa(volume),
			fmt.Sprintf("%.0f", marketCap),
			fmt.Sprintf("%.2f", peRatio),
			fmt.Sprintf("%.2f", dividendYield),
			fmt.Sprintf("%.3f", beta),
			fmt.Sprintf("%.2f", week52High),
			fmt.Sprintf("%.2f", week52Low),
			fmt.Sprintf("%.2f", changePercent),
		}

		if err := writer.Write(row); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing row %d: %v\n", i, err)
			os.Exit(1)
		}

		// Afficher la progression tous les 100k lignes
		if i%100000 == 0 {
			fmt.Printf("Progress: %d rows written (%.1f%%)\n", i, float64(i)/float64(rows)*100)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "Error flushing writer: %v\n", err)
		os.Exit(1)
	}

	duration := time.Since(start)
	fmt.Printf("\n✓ Successfully generated %d rows in %v\n", rows, duration)
	fmt.Printf("Output file: %s\n", output)

	// Afficher la taille du fichier
	if info, err := file.Stat(); err == nil {
		size := float64(info.Size())
		unit := "B"
		if size > 1024*1024*1024 {
			size /= 1024 * 1024 * 1024
			unit = "GB"
		} else if size > 1024*1024 {
			size /= 1024 * 1024
			unit = "MB"
		} else if size > 1024 {
			size /= 1024
			unit = "KB"
		}
		fmt.Printf("File size: %.2f %s\n", size, unit)
	}
}
