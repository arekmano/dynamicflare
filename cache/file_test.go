package cache_test

import (
	"os"
	"testing"

	"github.com/arekmano/dynamicflare/cache"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	// Test data
	value := "value"
	filePath := ".tmp"
	c := cache.NewFileCache(filePath)

	// Execute
	err := c.Write(value)

	// Verify
	require.NoError(t, err)

	// Execute
	result, err := c.Read()

	// Verify
	require.NoError(t, err)
	require.Equal(t, result, value)

	// Cleanup
	require.NoError(t, os.Remove(filePath))
}
