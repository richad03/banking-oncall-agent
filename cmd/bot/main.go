package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"banking-oncall-agent/internal/agent"
	"banking-oncall-agent/internal/config"
	"banking-oncall-agent/internal/knowledge"
	"banking-oncall-agent/internal/slack"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting Banking On-Call Agent in %s environment", cfg.Environment)

	// Initialize knowledge base
	log.Println("Loading knowledge base...")
	kb := knowledge.NewKnowledgeBase(cfg.KnowledgeBasePath)
	if err := kb.Load(); err != nil {
		log.Fatalf("Failed to load knowledge base: %v", err)
	}
	log.Println("Knowledge base loaded successfully")

	// Initialize agent
	log.Println("Initializing banking agent...")
	bankingAgent := agent.NewBankingAgent(cfg.AnthropicAPIKey, kb)
	log.Println("Banking agent initialized")

	// Initialize Slack bot
	log.Println("Initializing Slack bot...")
	bot, err := slack.NewBot(cfg.SlackBotToken, cfg.SlackAppToken, bankingAgent)
	if err != nil {
		log.Fatalf("Failed to create Slack bot: %v", err)
	}

	// Start bot in a goroutine
	go func() {
		log.Println("Starting Slack bot...")
		if err := bot.Start(); err != nil {
			log.Fatalf("Failed to start Slack bot: %v", err)
		}
	}()

	// Wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log.Println("Banking On-Call Agent is running... Press Ctrl+C to stop")
	<-interrupt

	log.Println("Shutting down Banking On-Call Agent...")
}
