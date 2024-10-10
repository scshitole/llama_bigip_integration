package processor

import (
	"encoding/json"
	"fmt"
	"strings"

	"llama_bigip_integration/bigip"
	"llama_bigip_integration/llama"
)

type Processor struct {
	llamaClient *llama.Client
	bigipClient *bigip.Client
}

func NewProcessor(llamaClient *llama.Client, bigipClient *bigip.Client) *Processor {
	return &Processor{
		llamaClient: llamaClient,
		bigipClient: bigipClient,
	}
}

func (p *Processor) ProcessQuery(query string) (string, error) {
	// First, get a completion from LLaMA
	llamaPrompt := fmt.Sprintf("Given the following query about F5 BIG-IP infrastructure, determine if it requires querying the BIG-IP system for virtual server information. If it does, respond with 'QUERY_BIGIP'. If not, provide a general answer. Query: %s", query)
	llamaResponse, err := p.llamaClient.GetCompletion(llamaPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to get completion from LLaMA: %v", err)
	}

	if strings.TrimSpace(llamaResponse) == "QUERY_BIGIP" {
		// Query BIG-IP for virtual server information
		virtualServers, err := p.bigipClient.GetVirtualServers()
		if err != nil {
			return "", fmt.Errorf("failed to get virtual servers from BIG-IP: %v", err)
		}

		// Convert virtual servers to JSON for easier processing by LLaMA
		vsJSON, err := json.Marshal(virtualServers)
		if err != nil {
			return "", fmt.Errorf("failed to marshal virtual servers to JSON: %v", err)
		}

		// Ask LLaMA to process the virtual server information and answer the query
		finalPrompt := fmt.Sprintf("Given the following query and virtual server information from F5 BIG-IP, provide a concise and informative answer. Query: %s\n\nVirtual Server Information: %s", query, string(vsJSON))
		finalResponse, err := p.llamaClient.GetCompletion(finalPrompt)
		if err != nil {
			return "", fmt.Errorf("failed to get final completion from LLaMA: %v", err)
		}

		return finalResponse, nil
	}

	// If BIG-IP querying is not required, return the LLaMA response directly
	return llamaResponse, nil
}
