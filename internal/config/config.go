package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	SlackBotToken     string
	SlackAppToken     string
	AnthropicAPIKey   string
	KnowledgeBasePath string
	Port              string
	Environment       string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	config := &Config{
		SlackBotToken:     os.Getenv("SLACK_BOT_TOKEN"),
		SlackAppToken:     os.Getenv("SLACK_APP_TOKEN"),
		AnthropicAPIKey:   os.Getenv("ANTHROPIC_API_KEY"),
		KnowledgeBasePath: getEnvOrDefault("KNOWLEDGE_BASE_PATH", "../banking-oncall"),
		Port:              getEnvOrDefault("PORT", "8080"),
		Environment:       getEnvOrDefault("ENVIRONMENT", "development"),
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// validate ensures required configuration is present
func (c *Config) validate() error {
	if c.SlackBotToken == "" {
		return fmt.Errorf("SLACK_BOT_TOKEN is required")
	}
	if c.SlackAppToken == "" {
		return fmt.Errorf("SLACK_APP_TOKEN is required")
	}
	if c.AnthropicAPIKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY is required")
	}
	return nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
