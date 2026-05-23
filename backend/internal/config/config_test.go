package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	// Test when environment variable is not set
	os.Unsetenv("TEST_KEY_NOT_SET")
	val := getEnvOrDefault("TEST_KEY_NOT_SET", "default_val")
	assert.Equal(t, "default_val", val)

	// Test when environment variable is set
	os.Setenv("TEST_KEY_SET", "custom_val")
	val = getEnvOrDefault("TEST_KEY_SET", "default_val")
	assert.Equal(t, "custom_val", val)
	
	// Cleanup
	os.Unsetenv("TEST_KEY_SET")
}

func TestGetAsynqRedisOpt(t *testing.T) {
	os.Setenv("REDIS_HOST", "127.0.0.1")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_PASSWORD", "secret")

	opt := GetAsynqRedisOpt()
	assert.Equal(t, "127.0.0.1:6380", opt.Addr)
	assert.Equal(t, "secret", opt.Password)

	// Cleanup
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("REDIS_PASSWORD")
}
