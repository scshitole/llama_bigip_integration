package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"llama_bigip_integration/bigip"
	"llama_bigip_integration/llama"
	"llama_bigip_integration/processor"
)

func main() {
	llamaURL := os.Getenv("LLAMA_API_URL")
	if llamaURL == "" {
		llamaURL = "http://localhost:5000/generate" // Default URL
	}
	llamaClient, err := llama.NewClient(llamaURL)
	if err != nil {
		log.Fatalf("Failed to initialize LLaMA client: %v", err)
	}

	var bigipClient *bigip.Client
	bigipHost := os.Getenv("BIGIP_HOST")
	bigipUsername := os.Getenv("BIGIP_USERNAME")
	bigipPassword := os.Getenv("BIGIP_PASSWORD")
	if bigipHost == "" || bigipUsername == "" || bigipPassword == "" {
		bigipClient = bigip.NewMockClient()
	} else {
		bigipClient, err = bigip.NewClient(bigipHost, bigipUsername, bigipPassword, false)
		if err != nil {
			log.Fatalf("Failed to initialize BIG-IP client: %v", err)
		}
	}

	queryProcessor := processor.NewProcessor(llamaClient, bigipClient)

	if len(os.Args) > 1 {
		query := strings.Join(os.Args[1:], " ")
		handleQuery(queryProcessor, bigipClient, query)
	} else {
		fmt.Println("Please provide a query as a command-line argument.")
	}
}

func handleQuery(queryProcessor *processor.Processor, bigipClient *bigip.Client, query string) {
	response, err := queryProcessor.ProcessQuery(query)
	if err != nil {
		rawData, err := getRawData(bigipClient, query)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println(rawData)
		}
	} else {
		fmt.Println(response)
	}
}

func getRawData(client *bigip.Client, query string) (string, error) {
	var data interface{}
	var err error

	queryLower := strings.ToLower(query)
	switch {
	case strings.Contains(queryLower, "virtual server"):
		data, err = client.GetVirtualServers()
	case strings.Contains(queryLower, "pool member"):
		pools, err := client.GetPools()
		if err != nil {
			return "", fmt.Errorf("failed to get pools: %v", err)
		}

		allPoolMembers := make(map[string][]bigip.PoolMember)
		for _, pool := range pools {
			members, err := client.GetPoolMembers(pool.Name)
			if err != nil {
				return "", fmt.Errorf("failed to get pool members for '%s': %v", pool.Name, err)
			}
			allPoolMembers[pool.Name] = members
		}
		data = allPoolMembers
	case strings.Contains(queryLower, "pool"):
		data, err = client.GetPools()
	default:
		return "", fmt.Errorf("unable to determine data type from query")
	}

	if err != nil {
		return "", fmt.Errorf("failed to get data: %v", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal data to JSON: %v", err)
	}

	return string(jsonData), nil
}
