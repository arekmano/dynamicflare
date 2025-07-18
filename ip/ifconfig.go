package ip

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// IfConfigClient is the client used to fetch your public IP from ifconfig.co
type IfConfigClient struct {
	logger *logrus.Entry
	Client *http.Client
}

// NewIfConfigClient creates a new client.
func NewIfConfigClient(logger *logrus.Entry) *IfConfigClient {
	return &IfConfigClient{
		Client: http.DefaultClient,
		logger: logger.WithField("component", "IfConfigClient"),
	}
}

// GetPublicIP returns the public IP
func (c *IfConfigClient) GetPublicIP() (newIP string, err error) {
	c.logger.Debug("Fetching Public IP")
	resp, err := c.Client.Get("https://ifconfig.me/ip")
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("error fetching public ip")
	}
	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	newIP = strings.Trim(string(ipBytes), "\n")
	c.logger.
		WithField("IP", newIP).
		Debug("Public IP Fetched")
	return
}
