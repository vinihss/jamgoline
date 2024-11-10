// File: jamgoline/cmd/jamgoline/main.go
package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vinihss/jamgoline/internal/config"
	"github.com/vinihss/jamgoline/pkg/agent"
	"github.com/vinihss/jamgoline/pkg/pubsub"
	"log"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "jamgoline",
		Short: "jamgoline is a PubSub system with agents and templates.",
	}

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available agent templates",
	Run: func(cmd *cobra.Command, args []string) {
		handleListTemplates()
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the configured agent pipeline",
	Run: func(cmd *cobra.Command, args []string) {
		handleRunPipeline()
	},
}

// handleListTemplates lists all registered templates in the system.
func handleListTemplates() {
	templates := agent.ListTemplates()
	fmt.Println("Available Templates:")
	for _, tmpl := range templates {
		fmt.Println("-", tmpl)
	}
}

// handleRunPipeline runs agents as configured in agents_config.yaml.
func handleRunPipeline() {
	conf, err := config.LoadConfig("config/agents_config.yaml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Define secret key for PubSub and initialize system
	secretKey := []byte("mysecretkey12345")
	pubsubSystem, err := pubsub.NewPubSub(secretKey)
	if err != nil {
		log.Fatalf("Error creating PubSub system: %v", err)
	}

	// Register and start agents based on configuration
	agents := make(map[string]*agent.Agent)
	for _, agentConfig := range conf.Agents {
		newAgent, err := agent.NewAgentFromTemplate(agentConfig.Template, agentConfig.Name, agentConfig.Topics, pubsubSystem)
		if err != nil {
			log.Fatalf("Error creating agent %s: %v", agentConfig.Name, err)
		}
		agents[agentConfig.Name] = newAgent
		go newAgent.Run()
	}

	// Execute pipeline sequence
	for _, agentName := range conf.Pipeline.Sequence {
		agent, exists := agents[agentName]
		if !exists {
			log.Fatalf("Agent %s not found in the pipeline", agentName)
		}
		pubsubSystem.PublishTopic(agentName, fmt.Sprintf("Message to %s", agentName))
	}

	// Keep program running to observe output
	select {}
}
