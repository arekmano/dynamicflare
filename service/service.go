package service

import (
	"errors"

	"github.com/arekmano/dynamicflare/cache"
	"github.com/arekmano/dynamicflare/dns"
	"github.com/arekmano/dynamicflare/ip"

	"github.com/sirupsen/logrus"
)

// DynamicFlare main service
type DynamicFlare struct {
	logger          *logrus.Entry
	DNSRecordClient *dns.CloudflareClient
	IPClient        *ip.IfConfigClient
	IPCache         *cache.FileCache
}

// Config the configuration used to create the service.
type Config struct {
	Cloudflare    dns.CloudflareConfig
	CacheFileName string
	Records       []dns.Record
}

// New create a new DynamicFlare
func New(config *Config, level logrus.Level) *DynamicFlare {
	logger := logrus.New()
	logger.SetLevel(level)

	return &DynamicFlare{
		DNSRecordClient: dns.NewCloudflareClient(config.Cloudflare, logrus.NewEntry(logger)),
		IPClient:        ip.NewIfConfigClient(logrus.NewEntry(logger)),
		logger:          logger.WithField("component", "DynamicFlare"),
		IPCache:         cache.NewFileCache(config.CacheFileName, logrus.NewEntry(logger)),
	}
}

// ListDomains lists all the domains associated with the account
func (s *DynamicFlare) ListDomains() error {
	zones, err := s.DNSRecordClient.Zones()
	if err != nil {
		return err
	}
	logrus.Info("Results")
	for i, zone := range zones {
		logrus.
			WithField("id", zone.ID).
			WithField("name", zone.Name).
			WithField("status", zone.Status).
			Info(i)
	}
	return nil
}

// ListDomainRecords lists all the records associated with the account
func (s *DynamicFlare) ListDomainRecords() error {
	zones, err := s.DNSRecordClient.Zones()
	if err != nil {
		return err
	}
	for _, zone := range zones {
		logrus.WithField("domain", zone.Name).Info("Records for domain")

		err = s.listDomainRecords(zone.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DynamicFlare) listDomainRecords(id string) error {
	records, err := s.DNSRecordClient.Records(id)
	if err != nil {
		return err
	}
	for i, record := range records {
		logrus.
			WithField("id", record.ID).
			WithField("name", record.Name).
			WithField("content", record.Content).
			WithField("type", record.Type).
			Info(i)
	}
	return nil
}

// Update updates the public IP of the given records
func (s *DynamicFlare) Update(dryRun bool, records []dns.Record) error {
	newIP, err := s.IPClient.GetPublicIP()
	if err != nil {
		return err
	}
	cacheContents, err := s.IPCache.Read()
	if err != nil {
		s.logger.
			WithField("cache", s.IPCache).
			Warn("could not read old IP from cache")
	}

	entry := s.logger.
		WithField("old-ip", cacheContents.IpAddress).
		WithField("cached-since", cacheContents.CacheTime.String()).
		WithField("new-ip", newIP)
	if cacheContents.IpAddress != newIP && dryRun {
		entry.Info("IP is different from cached. Not updating (dry-run on)")
		return nil
	} else if cacheContents.IpAddress != newIP {
		entry.Info("IP is different from cached. Updating Records")
		err = s.DNSRecordClient.UpdateMany(records, newIP)
		if err != nil {
			return err
		}
		return s.IPCache.Write(newIP)
	}
	entry.Info("IP is the same as the cached one")
	return nil
}

func (c *Config) Validate() error {
	if c.Cloudflare.Email == "" {
		return errors.New("configuration validation error: email not specified")
	}
	if c.Cloudflare.Key == "" {
		return errors.New("configuration validation error: cloudflare key not specified")
	}
	if c.CacheFileName == "" {
		return errors.New("configuration validation error: cache file name not specified")
	}
	for _, record := range c.Records {
		if record.ID == "" {
			return errors.New("configuration validation error: record ID not specified")
		}
		if record.Name == "" {
			return errors.New("configuration validation error: record name not specified")
		}
		if record.ZoneID == "" {
			return errors.New("configuration validation error: record zone ID not specified")
		}
		if record.Type == "" {
			return errors.New("configuration validation error: record type not specified")
		}
	}
	return nil
}
