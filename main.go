package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-tfe"
	"github.com/queone/utl"
)

const (
	prgname = "tfe"
	prgver  = "0.1.1"
)

// Prints usage
func printUsage() {
	fmt.Printf(prgname + " v" + prgver + "\n" +
		"Terraform Cloud/Enterprise CLI utility. See https://github.com/queone/tfe\n" +
		"Usage: " + prgname + " [options]\n" +
		"  -o [filter]              List orgs, filter option\n" +
		"  -m [filter]              List modules, filter option\n" +
		"  -w [filter]              List workspaces, filter option\n" +
		"  -?, -h, --help           Print this usage page\n")
	os.Exit(0)
}

// Sets up the client
func SetupClient() *tfe.Client {
	// Retrieve API token and organization name from environment variables
	tfeDomain := os.Getenv("TF_DOMAIN")
	token := os.Getenv("TF_TOKEN")
	orgName := os.Getenv("TF_ORG")

	// Check if the environment variables are set
	if token == "" || orgName == "" || tfeDomain == "" {
		log.Fatal("One or more required environment variables (TF_TOKEN, TF_ORG, TF_DOMAIN) are not set.")
	}

	// Set up a configuration with your API token
	config := &tfe.Config{
		Token:   token,
		Address: tfeDomain,
	}

	// Create a new TFE client
	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating TFE client: %v", err)
	}

	return client
}

// Lists organizations
func ListOrganizations(client *tfe.Client, filter string) {
	orgs, err := client.Organizations.List(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error listing organization: %v", err)
	}
	if orgs.Items != nil && len(orgs.Items) > 0 {
		filter = strings.ToLower(filter)
		for _, o := range orgs.Items {
			name := strings.ToLower(o.Name)
			if utl.SubString(name, filter) {
				fmt.Printf("%s\n", o.Name)
			}
		}
	}
}

// Lists workspaces
func ListWorkspaces(client *tfe.Client, orgName string, filter string) {
	workspaces, err := client.Workspaces.List(context.Background(), orgName, &tfe.WorkspaceListOptions{})
	if err != nil {
		log.Fatalf("Error listing workspaces for organization %s: %v", orgName, err)
	}
	if workspaces.Items != nil && len(workspaces.Items) > 0 {
		filter = strings.ToLower(filter)
		for _, ws := range workspaces.Items {
			name := strings.ToLower(ws.Name)
			if utl.SubString(name, filter) {
				fmt.Printf("%s\n", ws.Name)
			}
		}
	}
}

// Lists registered modules
func ListModules(client *tfe.Client, orgName string, filter string) {
	options := tfe.RegistryModuleListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
	}
	var allModules []*tfe.RegistryModule

	for {
		modules, err := client.RegistryModules.List(context.Background(), orgName, &options)
		if err != nil {
			log.Fatalf("Error listing modules for organization %s: %v", orgName, err)
		}
		allModules = append(allModules, modules.Items...)
		if modules.NextPage == 0 {
			break
		}
		options.PageNumber = modules.NextPage
	}
	if len(allModules) > 0 {
		filter = strings.ToLower(filter)
		for _, m := range allModules {
			name := strings.ToLower(m.Name)
			if utl.SubString(name, filter) {
				for _, v := range m.VersionStatuses {
					vVer := v.Version
					vStat := v.Status
					vErr := v.Error
					fmt.Printf("%-80s %-6s %-6s %s\n",
						"localterraform.com/"+m.Namespace+"/"+m.Name+"/"+m.Provider,
						vVer, vStat, vErr)
				}
			}
		}
	}
}

func main() {
	numberOfArguments := len(os.Args[1:]) // Not including the program itself
	if numberOfArguments < 1 || numberOfArguments > 2 {
		printUsage() // Don't accept less than 1 or more than 2 arguments
	}
	switch numberOfArguments {
	case 1: // Process 1-argument requests
		arg1 := os.Args[1]
		switch arg1 {
		case "-?", "-h", "--help":
			printUsage()
		}
		client := SetupClient()
		switch arg1 {
		case "-o":
			ListOrganizations(client, "")
		case "-m":
			ListModules(client, os.Getenv("TF_ORG"), "")
		case "-w":
			ListWorkspaces(client, os.Getenv("TF_ORG"), "")
		}
	case 2: // Process 2-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		client := SetupClient()
		switch arg1 {
		case "-o":
			ListOrganizations(client, arg2)
		case "-m":
			ListModules(client, os.Getenv("TF_ORG"), arg2)
		case "-w":
			ListWorkspaces(client, os.Getenv("TF_ORG"), arg2)
		}
	default:
		printUsage()
	}
}
