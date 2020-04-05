package service

import (
	"github.com/arekmano/dynamicflare/cache"
	"github.com/arekmano/dynamicflare/dns"
	"github.com/arekmano/dynamicflare/ip"

	"github.com/sirupsen/logrus"
)

// DynamicFlare main service
type DynamicFlare struct {
	logger     *logrus.Entry
	cloudflare *dns.CloudflareClient
	ifconfig   *ip.IfConfigClient
	filecache  *cache.FileCache
}

type Config struct {
	Cloudflare    *dns.CloudflareConfig
	CacheFileName string
	Records       []dns.Record
}

// New create a new DynamicFlare
func New(config *Config) *DynamicFlare {
	return &DynamicFlare{
		cloudflare: dns.NewCloudflareClient(config.Cloudflare),
		ifconfig:   ip.NewIfConfigClient(),
		logger:     logrus.WithField("component", "DynamicFlare"),
		filecache:  cache.NewFileCache(config.CacheFileName),
	}
}

func (s *DynamicFlare) ListDomains() error {
	zones, err := s.cloudflare.Zones()
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

func (s *DynamicFlare) ListDomainRecords() error {
	zones, err := s.cloudflare.Zones()
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
	domains, err := s.cloudflare.Domains(id)
	if err != nil {
		return err
	}
	for i, domain := range domains {
		logrus.
			WithField("id", domain.ID).
			WithField("name", domain.Name).
			WithField("content", domain.Content).
			WithField("type", domain.Type).
			Info(i)
	}
	return nil
}

// Run run the service
func (s *DynamicFlare) Run(dryRun bool, records []dns.Record) error {
	newIP, err := s.ifconfig.GetPublicIP()
	if err != nil {
		return err
	}
	ip, err := s.filecache.Read()
	if err != nil {
		return err
	}

	entry := s.logger.
		WithField("old-ip", ip).
		WithField("new-ip", newIP)
	if ip != newIP && dryRun {
		entry.Info("IP is different from cached. Not updating (dry-run on)")
		return nil
	} else if ip != newIP {
		entry.Info("IP is different from cached. Updating Records")
		return s.cloudflare.UpdateMany(records, newIP)
	}
	entry.Info("IP is the same as the cached one")
	return nil
}
