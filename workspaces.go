package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/git719/utl"
	"github.com/hashicorp/go-tfe"
)

func ListWorkspaces(client *tfe.Client, orgName string, filter string) {
	var allWorkspaces []*tfe.Workspace

	options := tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100, // Set the page size to a reasonable number
		},
	}

	// Retrieve all workspaces from the organization
	for {
		// Lists workspaces
		workspaces, err := client.Workspaces.List(context.Background(), orgName, &options)
		if err != nil {
			log.Fatalf("Error listing workspaces for organization %s: %v", orgName, err)
		}

		// Append the retrieved workspaces to the allWorkspaces slice
		allWorkspaces = append(allWorkspaces, workspaces.Items...)

		// Break the loop if there are no more pages
		if workspaces.Pagination.NextPage == 0 {
			break
		}

		// Set the next page for the next request
		options.PageNumber = workspaces.Pagination.NextPage
	}

	// Process the retrieved workspaces
	filter = strings.ToLower(filter)
	for _, ws := range allWorkspaces {
		name := strings.ToLower(ws.Name)
		if strings.Contains(name, filter) { // Assuming `utl.SubString` checks for substring
			fmt.Printf("%s\n", ws.Name)
		}
	}
}

func ShowWorkspace(client *tfe.Client, orgName string, wsName string) {
	// Shows details of a specific workspace
	workspace, err := client.Workspaces.Read(context.Background(), orgName, wsName)
	if err != nil {
		log.Fatalf("Error retrieving workspace %s in organization %s: %v", wsName, orgName, err)
	}

	colon := utl.Whi(":")

	desc := workspace.Description
	if desc == "" {
		desc = `""`
	}
	workingDir := workspace.WorkingDirectory
	if workingDir == "" {
		workingDir = `""`
	}

	fmt.Printf("%s%s %s\n", utl.Blu("workspace_name"), colon, utl.Gre(workspace.Name))
	fmt.Printf("%s%s %s\n", utl.Blu("workspace_id"), colon, utl.Gre(workspace.ID))
	fmt.Printf("%s%s %s\n", utl.Blu("created_at"), colon, utl.Gre(workspace.CreatedAt.Format("2006-01-02 15:04")))
	fmt.Printf("%s%s %s\n", utl.Blu("updated_at"), colon, utl.Gre(workspace.UpdatedAt.Format("2006-01-02 15:04")))
	fmt.Printf("%s%s %s\n", utl.Blu("description"), colon, utl.Gre(desc))
	fmt.Printf("%s%s %s\n", utl.Blu("terraform_version"), colon, utl.Gre(workspace.TerraformVersion))
	fmt.Printf("%s%s %s\n", utl.Blu("auto_apply"), colon, utl.Gre(workspace.AutoApply))
	fmt.Printf("%s%s %s\n", utl.Blu("working_directory"), colon, utl.Gre(workingDir))

	// Display Execution Mode
	fmt.Printf("%s%s %s\n", utl.Blu("execution_mode"), colon, utl.Gre(workspace.ExecutionMode))

	// Attempt to display Agent Pool ID if Execution Mode is "agent"
	if workspace.ExecutionMode == "agent" {
		if workspace.AgentPool != nil {
			// Assuming AgentPoolID is the correct field to access the ID string
			agentPoolID := workspace.AgentPool.ID
			agentPool, err := client.AgentPools.Read(context.Background(), agentPoolID)
			if err != nil {
				log.Fatalf("Error retrieving agent pool %s: %v", agentPoolID, err)
			}
			fmt.Printf("  %s%s %s\n", utl.Blu("agent_pool_name"), colon, utl.Gre(agentPool.Name))
		} else {
			fmt.Printf("  %s%s %s\n", utl.Blu("agent_pool_id"), colon, utl.Gre("Not available"))
		}
	}

	// Fetch and display environment variables and terraform variables
	variables, err := client.Variables.List(context.Background(), workspace.ID, &tfe.VariableListOptions{})
	if err != nil {
		log.Fatalf("Error retrieving variables for workspace %s: %v", wsName, err)
	}
	fmt.Printf("%s%s\n", utl.Blu("variables"), colon)

	// Initialize flags to check if there are any variables in each category
	hasEnvVars := false
	hasTerraformVars := false

	// First pass to check for variables in each category
	for _, variable := range variables.Items {
		if variable.Category == tfe.CategoryEnv {
			hasEnvVars = true
		} else if variable.Category == tfe.CategoryTerraform {
			hasTerraformVars = true
		}
	}

	// Print environment variables if they exist
	if hasEnvVars {
		fmt.Printf("  %s%s\n", utl.Blu("environment"), colon)
		for _, variable := range variables.Items {
			if variable.Category == tfe.CategoryEnv {
				fmt.Printf("    %s: %s\n", utl.Blu(variable.Key), utl.Gre(variable.Value))
			}
		}
	}

	// Print terraform variables if they exist
	if hasTerraformVars {
		fmt.Printf("  %s%s\n", utl.Blu("terraform"), colon)
		for _, variable := range variables.Items {
			if variable.Category == tfe.CategoryTerraform {
				fmt.Printf("    %s: %s\n", utl.Blu(variable.Key), utl.Gre(variable.Value))
			}
		}
	}
}

func CloneWorkspace(client *tfe.Client, orgName, srcWsName, destWsName string) {
	// Clones a workspace from WS_SRC to WS_DES, including variables

	// Get source workspace details
	srcWorkspace, err := client.Workspaces.Read(context.Background(), orgName, srcWsName)
	if err != nil {
		log.Fatalf("Error retrieving source workspace %s in organization %s: %v", srcWsName, orgName, err)
	}

	// Create new workspace with the same attributes as the source, but with a new name
	options := tfe.WorkspaceCreateOptions{
		Name:             tfe.String(destWsName),
		AutoApply:        tfe.Bool(srcWorkspace.AutoApply),
		TerraformVersion: tfe.String(srcWorkspace.TerraformVersion),
		WorkingDirectory: tfe.String(srcWorkspace.WorkingDirectory),
		Description:      tfe.String(srcWorkspace.Description),
		ExecutionMode:    tfe.String(srcWorkspace.ExecutionMode),
	}

	// If the source workspace uses an agent pool, set the AgentPoolID for the new workspace
	if srcWorkspace.ExecutionMode == "agent" && srcWorkspace.AgentPool != nil {
		options.AgentPoolID = tfe.String(srcWorkspace.AgentPool.ID)
	}

	destWorkspace, err := client.Workspaces.Create(context.Background(), orgName, options)
	if err != nil {
		log.Fatalf("Error creating destination workspace %s: %v", utl.Gre(destWsName), err)
	}

	// Fetch variables from the source workspace
	variables, err := client.Variables.List(context.Background(), srcWorkspace.ID, &tfe.VariableListOptions{})
	if err != nil {
		log.Fatalf("Error retrieving variables for source workspace %s: %v", utl.Blu(srcWsName), err)
	}

	// Copy variables to the destination workspace
	for _, variable := range variables.Items {
		createVariableOptions := tfe.VariableCreateOptions{
			Key:       tfe.String(variable.Key),
			Value:     tfe.String(variable.Value),
			Category:  tfe.Category(variable.Category),
			HCL:       tfe.Bool(variable.HCL),
			Sensitive: tfe.Bool(variable.Sensitive),
		}
		_, err := client.Variables.Create(context.Background(), destWorkspace.ID, createVariableOptions)
		if err != nil {
			log.Fatalf("Error creating variable %s in destination workspace %s: %v", variable.Key, destWsName, err)
		}
	}

	fmt.Printf("Successfully cloned workspace %s to %s\n", utl.Blu(srcWsName), utl.Gre(destWsName))
}
