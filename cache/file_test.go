package cache_test

import (
	"os"
	"testing"
	"time"

	"github.com/arekmano/dynamicflare/cache"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	// Test data
	value := "value"
	filePath := ".tmp"
	c := cache.NewFileCache(filePath, logrus.WithTime(time.Now()))

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
