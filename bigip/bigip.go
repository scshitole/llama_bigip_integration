package bigip

import (
	"crypto/tls"
	"fmt"

	"github.com/f5devcentral/go-bigip"
)

type Client struct {
	bigip *bigip.BigIP
}

func NewClient(host, username, password string, sslVerify bool) (*Client, error) {
	config := &bigip.Config{
		Address:  host,
		Username: username,
		Password: password,
	}
	bigipClient := bigip.NewSession(config)
	bigipClient.Transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: !sslVerify}
	return &Client{bigip: bigipClient}, nil
}

func (c *Client) GetVirtualServers() ([]bigip.VirtualServer, error) {
	virtualServers, err := c.bigip.VirtualServers()
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual servers: %v", err)
	}

	return virtualServers.VirtualServers, nil
}

func (c *Client) GetVirtualServerByName(name string) (*bigip.VirtualServer, error) {
	virtualServer, err := c.bigip.GetVirtualServer(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual server '%s': %v", name, err)
	}

	return virtualServer, nil
}
