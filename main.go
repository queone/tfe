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
	prgver  = "0.3.0"
)

// Prints usage
func printUsage() {
	fmt.Printf(prgname + " v" + prgver + "\n" +
		"Terraform Cloud/Enterprise CLI utility. See https://github.com/queone/tfe\n" +
		"Usage: " + prgname + " [options]\n" +
		"  -o [filter]              List orgs, filter option\n" +
		"  -m [filter]              List most recent version of modules, filter option\n" +
		"  -ma [filter]             List all version of modules, filter option\n" +
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

// Lists organizations, with a name filter option
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

// Lists workspaces, with a name filter option
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

// Lists registered modules, with a name filter option, which version option
func ListModules(client *tfe.Client, orgName string, filter string, ver string) {
	options := tfe.RegistryModuleListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
	}
	var allModules []*tfe.RegistryModule

	// Retrieve all modules from the organization
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

	if len(allModules) == 1 {
		modID := tfe.NewPrivateRegistryModuleID(orgName, allModules[0].Name, allModules[0].Provider)
		PrintModuleDetails(client, modID)
	} else if len(allModules) > 1 {
		filter = strings.ToLower(filter)

		if ver == "all" {
			// Print all versions of each filtered module
			for _, m := range allModules {
				if utl.SubString(strings.ToLower(m.Name), filter) {
					for _, v := range m.VersionStatuses {
						fmt.Printf("%-80s %-10s %s\n", "localterraform.com/"+m.Namespace+"/"+m.Name+"/"+m.Provider, v.Version, v.Status)
					}
				}
			}
		} else {
			// Map to store the latest version of each module
			latestVersions := make(map[string]int)

			// Iterate over all modules and track the latest version
			for _, m := range allModules {
				name := strings.ToLower(m.Name)
				if utl.SubString(name, filter) {
					for i, v := range m.VersionStatuses {
						if current, exists := latestVersions[m.Name]; !exists || strings.Compare(v.Version, m.VersionStatuses[current].Version) > 0 {
							latestVersions[m.Name] = i
						}
					}
				}
			}

			// Print the latest version of each filtered module
			for _, m := range allModules {
				// if v, exists := latestVersions[m.Name]; exists {
				// 	vVer := m.VersionStatuses[latestVersions[m.Name]].Version
				// 	vStat := m.VersionStatuses[latestVersions[m.Name]].Status
				// 	fmt.Printf("%-80s %-10s %s\n", "localterraform.com/"+m.Namespace+"/"+m.Name+"/"+m.Provider, vVer, vStat)
				// }
				if i, exists := latestVersions[m.Name]; exists {
					v := m.VersionStatuses[i]
					fmt.Printf("%-80s %-10s %s\n", "localterraform.com/"+m.Namespace+"/"+m.Name+"/"+m.Provider, v.Version, v.Status)
				}
			}
		}
	}
}

// Prints details of the module
func PrintModuleDetails(client *tfe.Client, modID tfe.RegistryModuleID) {
	module, err := client.RegistryModules.Read(context.Background(), modID)
	if err != nil {
		log.Fatalf("Error retrieving module with ID %v: %v", modID, err)
	}
	fmt.Printf("%+v\n", *module)
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
			ListModules(client, os.Getenv("TF_ORG"), "", "latest")
		case "-ma":
			ListModules(client, os.Getenv("TF_ORG"), "", "all")
		case "-w":
			ListWorkspaces(client, os.Getenv("TF_ORG"), "")
		}
	case 2: // Process 2-argument requests
		arg1 := os.Args[1]
		filter := os.Args[2]
		client := SetupClient()
		switch arg1 {
		case "-o":
			ListOrganizations(client, filter)
		case "-m":
			ListModules(client, os.Getenv("TF_ORG"), filter, "latest")
		case "-ma":
			ListModules(client, os.Getenv("TF_ORG"), filter, "all")
		case "-w":
			ListWorkspaces(client, os.Getenv("TF_ORG"), filter)
		}
	default:
		printUsage()
	}
}
