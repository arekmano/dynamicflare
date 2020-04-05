package dns

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CloudflareClient struct {
	logger *logrus.Entry
	Client *http.Client
	Key    string
	Email  string
}

type UpdateRecordResponse struct {
	Result   Record                   `json:"result"`
	Success  bool                     `json:"success"`
	Errors   []map[string]interface{} `json:"errors"`
	Messages []map[string]interface{} `json:"messages"`
}

type ListRecordsResponse struct {
	Result   []Record                 `json:"result"`
	Success  bool                     `json:"success"`
	Errors   []map[string]interface{} `json:"errors"`
	Messages []map[string]interface{} `json:"messages"`
}

type ListDomainsResponse struct {
	Result   []Domain                 `json:"result"`
	Success  bool                     `json:"success"`
	Errors   []map[string]interface{} `json:"errors"`
	Messages []map[string]interface{} `json:"messages"`
}

type Domain struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Record struct {
	ID         string `json:"id"`
	RecordType string `json:"type"`
	Name       string `json:"name"`
	Content    string `json:"content"`
}

type CloudflareConfig struct {
	Key   string
	Email string
}

func NewCloudflareClient(config *CloudflareConfig) *CloudflareClient {
	return &CloudflareClient{
		Email:  config.Email,
		Key:    config.Key,
		Client: http.DefaultClient,
		logger: logrus.WithField("component", "CloudflareClient"),
	}
}

func (c *CloudflareClient) UpdateMany(records []Record, newIP string) error {
	for _, record := range records {
		err := c.Update(record, newIP)
		if err != nil {
			return errors.Wrapf(err, "Error updating %s", record.Name)
		}
	}
	return nil
}

func (c *CloudflareClient) Domains(id string) ([]Record, error) {
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones/"+id+"/dns_records", strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Key", c.Key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Email", c.Email)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response ListRecordsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, errors.New("Incorrect status returned")
	}
	return response.Result, err
}

func (c *CloudflareClient) Zones() ([]Domain, error) {
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones", strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Key", c.Key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Email", c.Email)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response ListDomainsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, errors.New("Incorrect status returned")
	}
	return response.Result, err
}

func (c *CloudflareClient) send(record Record) (*http.Response, error) {
	bodyBytes, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", "https://api.cloudflare.com/client/v4/zones/2a1c86b70e031a2d3f1a45ce2bbfa544/dns_records/"+record.ID, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Key", c.Key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Email", c.Email)

	c.logger.
		WithField("name", record.Name).
		WithField("type", record.RecordType).
		WithField("new-ip", record.Content).
		Debug("sending request to Cloudflare")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *CloudflareClient) Update(record Record, newIP string) error {
	record.Content = newIP
	resp, err := c.send(record)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	responseRecord := UpdateRecordResponse{}
	err = json.Unmarshal(body, &responseRecord)
	if err != nil {
		return err
	}
	if responseRecord.Success {
		c.logger.
			WithField("name", record.Name).
			Info("Cloudflare DNS record successfully updated")
	} else {
		c.logger.
			WithField("name", record.Name).
			WithField("response", responseRecord).
			WithField("status-code", resp.StatusCode).
			WithField("status", resp.Status).
			Error("Error Updating Cloudflare DNS record")
	}
	return nil
}
