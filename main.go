package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"llama_bigip_integration/bigip"
	"llama_bigip_integration/llama"
	"llama_bigip_integration/processor"
)

func main() {
	// Check if a query was provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go \"<your query>\"")
		os.Exit(1)
	}

	// Join all arguments after the program name as the query
	query := strings.Join(os.Args[1:], " ")

	// Initialize LLaMA client
	llamaURL := os.Getenv("LLAMA_API_URL")
	if llamaURL == "" {
		llamaURL = "http://localhost:5000/generate" // Default URL
	}
	llamaClient, err := llama.NewClient(llamaURL)
	if err != nil {
		log.Fatalf("Failed to initialize LLaMA client: %v", err)
	}

	// Initialize BIG-IP client
	bigipHost := os.Getenv("BIGIP_HOST")
	bigipUsername := os.Getenv("BIGIP_USERNAME")
	bigipPassword := os.Getenv("BIGIP_PASSWORD")
	if bigipHost == "" || bigipUsername == "" || bigipPassword == "" {
		log.Fatalf("BIG-IP environment variables not set. Please set BIGIP_HOST, BIGIP_USERNAME, and BIGIP_PASSWORD")
	}
	bigipClient, err := bigip.NewClient(bigipHost, bigipUsername, bigipPassword, false)
	if err != nil {
		log.Fatalf("Failed to initialize BIG-IP client: %v", err)
	}

	// Initialize query processor
	queryProcessor := processor.NewProcessor(llamaClient, bigipClient)

	// Process the query
	response, err := queryProcessor.ProcessQuery(query)
	if err != nil {
		log.Fatalf("Error processing query: %v", err)
	}

	fmt.Println("Response:", response)
}
