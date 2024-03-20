package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"
)

const (
	prgname = "tfe"
	prgver  = "0.0.1"
)

func printUsage() {
	fmt.Printf(prgname + " v" + prgver + "\n" +
		"Terraform Cloud API utility. See https://github.com/queone/tfe\n" +
		"Usage: " + prgname + " [options]\n" +
		"  -?, -h, --help                    Print this usage page\n")
	os.Exit(0)
}

func setupClient() {
	// Retrieve API token and organization name from environment variables
	token := os.Getenv("TF_TOKEN")
	orgName := os.Getenv("TF_ORG")

	// Check if the environment variables are set
	if token == "" {
		log.Fatal("TF_TOKEN is not set.")
	}
	if orgName == "" {
		log.Fatal("TF_ORG is not set.")
	}

	// Set up a configuration with your API token
	config := &tfe.Config{
		Token: token,
	}

	// Create a new TFE client
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating TFE client: %v", err)
	}

	// Use the client to list workspaces in an organization
	workspaces, err := client.Workspaces.List(context.Background(), orgName, &tfe.WorkspaceListOptions{})
	if err != nil {
		log.Fatalf("Error listing workspaces for organization %s: %v", orgName, err)
	}

	for _, ws := range workspaces.Items {
		fmt.Printf("Workspace: %s\n", ws.Name)
	}
}

func main() {
	numberOfArguments := len(os.Args[1:]) // Not including the program itself
	if numberOfArguments < 1 || numberOfArguments > 4 {
		printUsage() // Don't accept less than 1 or more than 4 arguments
	}
	switch numberOfArguments {
	case 1: // Process 1-argument requests
		arg1 := os.Args[1]
		// This first set of 1-arg requests do not require API tokens to be set up
		switch arg1 {
		case "-?", "-h", "--help":
			printUsage()
		case "1":
			setupClient()
		}
	default:
		printUsage()
	}
}
