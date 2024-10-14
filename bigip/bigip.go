package bigip

import (
	"crypto/tls"
	"fmt"

	"github.com/f5devcentral/go-bigip"
)

type Client struct {
	bigip *bigip.BigIP
	mock  bool
}

// PoolMember represents a member of a pool
type PoolMember struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	State   string `json:"state"`
}

func NewClient(host, username, password string, sslVerify bool) (*Client, error) {
	config := &bigip.Config{
		Address:  host,
		Username: username,
		Password: password,
	}
	bigipClient := bigip.NewSession(config)
	bigipClient.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: !sslVerify}
	return &Client{bigip: bigipClient, mock: false}, nil
}

func NewMockClient() *Client {
	return &Client{mock: true}
}

func (c *Client) GetVirtualServers() ([]bigip.VirtualServer, error) {
	if c.mock {
		return []bigip.VirtualServer{
			{Name: "vs1", Destination: "10.0.0.1:80"},
			{Name: "vs2", Destination: "10.0.0.2:443"},
		}, nil
	}
	virtualServers, err := c.bigip.VirtualServers()
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual servers: %v", err)
	}
	return virtualServers.VirtualServers, nil
}

func (c *Client) GetVirtualServerByName(name string) (*bigip.VirtualServer, error) {
	if c.mock {
		return &bigip.VirtualServer{Name: name, Destination: "10.0.0.1:80"}, nil
	}
	virtualServer, err := c.bigip.GetVirtualServer(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual server '%s': %v", name, err)
	}
	return virtualServer, nil
}

func (c *Client) GetPools() ([]bigip.Pool, error) {
	if c.mock {
		return []bigip.Pool{
			{Name: "pool1", Monitor: "/Common/http"},
			{Name: "pool2", Monitor: "/Common/https"},
		}, nil
	}
	pools, err := c.bigip.Pools()
	if err != nil {
		return nil, fmt.Errorf("failed to get pools: %v", err)
	}
	return pools.Pools, nil
}

func (c *Client) GetPoolByName(name string) (*bigip.Pool, error) {
	if c.mock {
		return &bigip.Pool{Name: name, Monitor: "/Common/http"}, nil
	}
	pool, err := c.bigip.GetPool(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool '%s': %v", name, err)
	}
	return pool, nil
}

func (c *Client) GetPoolMembers(poolName string) ([]PoolMember, error) {
	if c.mock {
		return []PoolMember{
			{Name: "member1", Address: "192.168.1.1", State: "up"},
			{Name: "member2", Address: "192.168.1.2", State: "down"},
		}, nil
	}
	members, err := c.bigip.PoolMembers(poolName)
	if err != nil {
		return nil, fmt.Errorf("failed to get pool members for '%s': %v", poolName, err)
	}
	var poolMembers []PoolMember
	for _, member := range members.PoolMembers {
		poolMembers = append(poolMembers, PoolMember{
			Name:    member.Name,
			Address: member.Address,
			State:   member.State,
		})
	}
	return poolMembers, nil
}
