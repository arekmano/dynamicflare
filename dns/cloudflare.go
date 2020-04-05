package dns

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	ZoneID  string `json:"zone_id"`
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
		c.logger.
			WithField("id", record.ID).
			WithField("name", record.Name).
			WithField("type", record.Type).
			WithField("zone_id", record.ZoneID).
			Debug("Updating record")
		err := c.Update(record, newIP)
		if err != nil {
			return errors.Wrapf(err, "Error updating %s", record.Name)
		}
	}
	return nil
}

func (c *CloudflareClient) Domains(id string) ([]Record, error) {
	body, err := c.send("GET", "https://api.cloudflare.com/client/v4/zones/"+id+"/dns_records", []byte(""))
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
	body, err := c.send("GET", "https://api.cloudflare.com/client/v4/zones", []byte(""))
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
func (c *CloudflareClient) send(verb, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(verb, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Key", c.Key)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Email", c.Email)

	c.logger.
		WithField("url", url).
		WithField("verb", verb).
		WithField("body", string(body)).
		Debug("sending request to Cloudflare")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.logger.
		WithField("status-code", resp.StatusCode).
		WithField("body", string(responseBody)).
		Debug("received response from Cloudflare")

	return responseBody, nil
}

func (c *CloudflareClient) Update(record Record, newIP string) error {
	record.Content = newIP
	bodyBytes, err := json.Marshal(record)
	if err != nil {
		return err
	}
	body, err := c.send("PUT", "https://api.cloudflare.com/client/v4/zones/"+record.ZoneID+"/dns_records/"+record.ID, bodyBytes)
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
			Error("Error Updating Cloudflare DNS record")
	}
	return nil
}
