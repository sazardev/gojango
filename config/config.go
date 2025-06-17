package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds application configuration
type Config struct {
	DatabaseURL string
	Debug       bool
	Port        string
	Host        string
	settings    map[string]interface{}
}

// New creates a new configuration with defaults
func New() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		Debug:       getEnvBool("DEBUG", false),
		Port:        getEnv("PORT", "8000"),
		Host:        getEnv("HOST", "localhost"),
		settings:    make(map[string]interface{}),
	}
}

// Set sets a configuration value
func (c *Config) Set(key string, value interface{}) {
	if c.settings == nil {
		c.settings = make(map[string]interface{})
	}
	c.settings[key] = value
}

// Get gets a configuration value with default
func (c *Config) Get(key string, defaultValue interface{}) interface{} {
	if val, exists := c.settings[key]; exists {
		return val
	}
	return defaultValue
}

// GetString gets a string configuration value
func (c *Config) GetString(key, defaultValue string) string {
	if val := c.Get(key, defaultValue); val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetInt gets an integer configuration value
func (c *Config) GetInt(key string, defaultValue int) int {
	if val := c.Get(key, defaultValue); val != nil {
		switch v := val.(type) {
		case int:
			return v
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return defaultValue
}

// GetBool gets a boolean configuration value
func (c *Config) GetBool(key string, defaultValue bool) bool {
	if val := c.Get(key, defaultValue); val != nil {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			return strings.ToLower(v) == "true" || v == "1"
		}
	}
	return defaultValue
}

// LoadFromEnv loads configuration from environment variables
func (c *Config) LoadFromEnv(prefix string) {
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := parts[0]
		value := parts[1]
		
		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}
		
		// Remove prefix and convert to lowercase with dots
		if prefix != "" {
			key = strings.TrimPrefix(key, prefix)
		}
		key = strings.ToLower(strings.ReplaceAll(key, "_", "."))
		
		c.Set(key, value)
	}
}

// getEnv gets environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool gets environment variable as boolean
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true" || value == "1"
	}
	return defaultValue
}
