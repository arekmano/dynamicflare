package service_test

import (
	"testing"

	"github.com/arekmano/dynamicflare/dns"
	"github.com/arekmano/dynamicflare/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func randomConfig() *service.Config {
	return &service.Config{
		Cloudflare: dns.CloudflareConfig{
			Email: uuid.New().String(),
			Key:   uuid.New().String(),
		},
		CacheFileName: uuid.New().String(),
		Records: []dns.Record{
			dns.Record{
				Content: uuid.New().String(),
				ID:      uuid.New().String(),
				Name:    uuid.New().String(),
				Type:    uuid.New().String(),
				ZoneID:  uuid.New().String(),
			},
		},
	}
}

func TestValidate(t *testing.T) {
	// Test Data
	testcases := [][]interface{}{
		{"Valid Data", randomConfig(), false},
	}

	noId := randomConfig()
	noId.Records[0].ID = ""
	testcases = append(testcases, []interface{}{"Invalid Data - no id", noId, true})

	noZone := randomConfig()
	noZone.Records[0].ZoneID = ""
	testcases = append(testcases, []interface{}{"Invalid Data - no zone id", noZone, true})

	noType := randomConfig()
	noType.Records[0].Type = ""
	testcases = append(testcases, []interface{}{"Invalid Data - no type", noType, true})

	noName := randomConfig()
	noName.Records[0].Name = ""
	testcases = append(testcases, []interface{}{"Invalid Data - no name", noName, true})

	noKey := randomConfig()
	noKey.Cloudflare.Key = ""
	testcases = append(testcases, []interface{}{"Invalid Data - no key", noKey, true})

	noEmail := randomConfig()
	noEmail.Cloudflare.Email = ""
	testcases = append(testcases, []interface{}{"Invalid Data - no email", noEmail, true})

	noCacheFile := randomConfig()
	noCacheFile.CacheFileName = ""
	testcases = append(testcases, []interface{}{"Invalid Data - no cache filename", noCacheFile, true})

	for _, testcase := range testcases {
		t.Run(testcase[0].(string), func(t *testing.T) {
			result := testcase[1].(*service.Config).Validate()

			if !testcase[2].(bool) {
				require.Nil(t, result)
			} else {
				require.Error(t, result)
			}
		})
	}
}
