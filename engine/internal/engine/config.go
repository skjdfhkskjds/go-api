package engine

import (
	"fmt"
	"os"
)

// Config represents the minimal application configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
}

// ServerConfig contains basic HTTP server configuration
type ServerConfig struct {
	Port         int `yaml:"port"`
	ReadTimeout  int `yaml:"read_timeout"`  // seconds
	WriteTimeout int `yaml:"write_timeout"` // seconds
	IdleTimeout  int `yaml:"idle_timeout"`  // seconds
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  10,
			WriteTimeout: 10,
			IdleTimeout:  60,
		},
	}
}

// LoadConfig loads configuration from a YAML file (placeholder for now)
// TODO: Implement YAML loading when dependencies are added
func LoadConfig(filename string) (*Config, error) {
	config := DefaultConfig()

	if filename == "" {
		return config, nil
	}

	// For now, just check if file exists and return default config
	// YAML parsing will be implemented later with dependencies
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return config, nil
	}

	return config, nil
}

// LoadConfigWithEnv loads configuration from YAML file and overrides with environment variables
func LoadConfigWithEnv(filename string) (*Config, error) {
	config, err := LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	// Override with environment variables
	// Basic port override for minimal config
	if port := os.Getenv("SERVER_PORT"); port != "" {
		// Simple conversion - in production, add proper error handling
		fmt.Sscanf(port, "%d", &config.Server.Port)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Basic validation rules
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}

	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}

	return nil
}
