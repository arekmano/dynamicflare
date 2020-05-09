package ip_test

import (
	"net"
	"testing"

	"github.com/arekmano/dynamicflare/ip"
	"github.com/stretchr/testify/require"
)

func TestGetPublicIP(t *testing.T) {
	// Test data
	client := ip.NewIfConfigClient()

	// Execute
	result, err := client.GetPublicIP()
	require.NoError(t, err)
	resultIP := net.ParseIP(result)
	require.NoError(t, err)

	// Execute
	result, err = client.GetPublicIP()
	require.NoError(t, err)
	resultIP2 := net.ParseIP(result)
	require.NoError(t, err)

	// Verify
	require.Equal(t, resultIP, resultIP2)
}
