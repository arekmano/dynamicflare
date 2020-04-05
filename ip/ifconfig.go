package ip

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type IfConfigClient struct {
	logger *logrus.Entry
	Client *http.Client
}

func NewIfConfigClient() *IfConfigClient {
	return &IfConfigClient{
		Client: http.DefaultClient,
		logger: logrus.WithField("component", "IfConfigClient"),
	}
}

func (c *IfConfigClient) GetPublicIP() (newIP string, err error) {
	c.logger.Debug("Fetching Public IP")
	resp, err := c.Client.Get("https://ifconfig.co/ip")
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("Error Fetching Public IP")
	}
	ipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	newIP = strings.Trim(string(ipBytes), "\n")
	c.logger.
		WithField("IP", newIP).
		Debug("Public IP Fetched")
	return
}
