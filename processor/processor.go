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
	llamaPrompt := fmt.Sprintf("Given the following query about F5 BIG-IP infrastructure, determine if it requires querying the BIG-IP system for virtual server, pool, or pool member information. If it does, respond with 'QUERY_BIGIP_VS' for virtual servers, 'QUERY_BIGIP_POOL' for pools, or 'QUERY_BIGIP_POOL_MEMBERS' for pool members. If not, provide a general answer. Query: %s", query)
	llamaResponse, err := p.llamaClient.GetCompletion(llamaPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to get completion from LLaMA: %v", err)
	}

	llamaResponse = strings.TrimSpace(llamaResponse)

	switch llamaResponse {
	case "QUERY_BIGIP_VS":
		return p.handleVirtualServerQuery(query)
	case "QUERY_BIGIP_POOL":
		return p.handlePoolQuery(query)
	case "QUERY_BIGIP_POOL_MEMBERS":
		return p.handlePoolMemberQuery(query)
	default:
		return llamaResponse, nil
	}
}

func (p *Processor) handleVirtualServerQuery(query string) (string, error) {
	virtualServers, err := p.bigipClient.GetVirtualServers()
	if err != nil {
		return "", fmt.Errorf("failed to get virtual servers from BIG-IP: %v", err)
	}

	var vsInfo []map[string]interface{}
	for _, vs := range virtualServers {
		vsData := map[string]interface{}{
			"name":        vs.Name,
			"destination": vs.Destination,
			"pool":        vs.Pool,
			"enabled":     vs.Enabled,
		}
		vsInfo = append(vsInfo, vsData)
	}

	vsJSON, err := json.Marshal(vsInfo)
	if err != nil {
		return "", fmt.Errorf("failed to marshal virtual servers to JSON: %v", err)
	}

	finalPrompt := fmt.Sprintf("Given the following query and virtual server information from F5 BIG-IP, provide a concise and informative answer. Include relevant details such as names, destinations, associated pools, and enabled status. Query: %s\n\nVirtual Server Information: %s", query, string(vsJSON))
	finalResponse, err := p.llamaClient.GetCompletion(finalPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to get final completion from LLaMA: %v", err)
	}

	return finalResponse, nil
}

func (p *Processor) handlePoolQuery(query string) (string, error) {
	pools, err := p.bigipClient.GetPools()
	if err != nil {
		return "", fmt.Errorf("failed to get pools from BIG-IP: %v", err)
	}

	poolJSON, err := json.Marshal(pools)
	if err != nil {
		return "", fmt.Errorf("failed to marshal pools to JSON: %v", err)
	}

	finalPrompt := fmt.Sprintf("Given the following query and pool information from F5 BIG-IP, provide a concise and informative answer. Query: %s\n\nPool Information: %s", query, string(poolJSON))
	finalResponse, err := p.llamaClient.GetCompletion(finalPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to get final completion from LLaMA: %v", err)
	}

	return finalResponse, nil
}

func (p *Processor) handlePoolMemberQuery(query string) (string, error) {
	pools, err := p.bigipClient.GetPools()
	if err != nil {
		return "", fmt.Errorf("failed to get pools from BIG-IP: %v", err)
	}

	var allPoolMembers []map[string]interface{}

	for _, pool := range pools {
		members, err := p.bigipClient.GetPoolMembers(pool.Name)
		if err != nil {
			return "", fmt.Errorf("failed to get pool members for '%s': %v", pool.Name, err)
		}

		poolInfo := map[string]interface{}{
			"poolName": pool.Name,
			"members":  members,
		}
		allPoolMembers = append(allPoolMembers, poolInfo)
	}

	poolMembersJSON, err := json.Marshal(allPoolMembers)
	if err != nil {
		return "", fmt.Errorf("failed to marshal pool members to JSON: %v", err)
	}

	finalPrompt := fmt.Sprintf("Given the following query and pool member information from F5 BIG-IP, provide a concise and informative answer. Query: %s\n\nPool Member Information: %s", query, string(poolMembersJSON))
	finalResponse, err := p.llamaClient.GetCompletion(finalPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to get final completion from LLaMA: %v", err)
	}

	return finalResponse, nil
}
