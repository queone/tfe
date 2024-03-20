package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/hashicorp/go-tfe"
)

func main() {
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
    workspaces, err := client.Workspaces.List(context.Background(), orgName, tfe.WorkspaceListOptions{})
    if err != nil {
        log.Fatalf("Error listing workspaces for organization %s: %v", orgName, err)
    }

    for _, ws := range workspaces.Items {
        fmt.Printf("Workspace: %s\n", ws.Name)
    }
}
