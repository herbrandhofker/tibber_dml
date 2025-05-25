package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/herbrandhofker/tibber_ddl/pkg/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type TokenEvent struct {
	TokenID  string `json:"token_id"`
	IsActive bool   `json:"is_active"`
	Action   string `json:"action"`
}

func printJSON(label string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error encoding %s: %v", label, err)
		return
	}
	fmt.Printf("\n%s:\n%s\n", label, string(jsonData))
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database connection string from environment
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Create database connection
	db, err := database.NewDatabase(connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Set up channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Listen for token changes
	tokenChan := make(chan string)
	go listenForTokenChanges(db, tokenChan)

	<-sigChan
	log.Println("Shutting down...")
}

func listenForTokenChanges(db database.DBInterface, tokenChan chan string) {
	// Listen for notifications on the token_changes channel
	_, err := db.Exec("LISTEN token_changes")
	if err != nil {
		log.Fatalf("Failed to listen for token changes: %v", err)
	}

	log.Println("Successfully subscribed to token_changes channel")

	// Wait for notification
	_, err = db.Exec("SELECT pg_notify('token_changes', '')")
	if err != nil {
		log.Printf("Error waiting for notification: %v", err)
		return
	}

	log.Println("Received notification on token_changes channel")

	// Process the event
	processTokenEvent(db, tokenChan)
}

func processTokenEvent(db database.DBInterface, tokenChan chan string) {
	// Query the latest token from the tibber_api_token table
	var token string
	err := db.QueryRow("SELECT token FROM tibber.tibber_api_token WHERE is_active = true ORDER BY created_at DESC LIMIT 1").Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No active token found")
			return
		}
		log.Printf("Error querying token: %v", err)
		return
	}

	if token != "" {
		log.Printf("Found active token, sending to channel")
		tokenChan <- token
	} else {
		log.Println("Token is not active, skipping data fetch")
	}
}

func fetchAndPrintData(db database.DBInterface, homeID string) {
	// Get consumption data
	hourlyConsumption, err := db.GetRecentHourlyConsumption(homeID)
	if err != nil {
		log.Printf("Error fetching hourly consumption: %v", err)
	} else {
		printJSON("Hourly Consumption", hourlyConsumption)
	}

	dailyConsumption, err := db.GetRecentDailyConsumption(homeID)
	if err != nil {
		log.Printf("Error fetching daily consumption: %v", err)
	} else {
		printJSON("Daily Consumption", dailyConsumption)
	}

	// Get production data
	hourlyProduction, err := db.GetRecentHourlyProduction(homeID)
	if err != nil {
		log.Printf("Error fetching hourly production: %v", err)
	} else {
		printJSON("Hourly Production", hourlyProduction)
	}

	dailyProduction, err := db.GetRecentDailyProduction(homeID)
	if err != nil {
		log.Printf("Error fetching daily production: %v", err)
	} else {
		printJSON("Daily Production", dailyProduction)
	}

	// Get price data
	prices, err := db.GetRecentPrices(homeID)
	if err != nil {
		log.Printf("Error fetching prices: %v", err)
	} else {
		printJSON("Prices", prices)
	}
}
