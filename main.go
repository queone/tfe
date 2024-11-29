package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"
)

const (
	prgname = "tfe"
	prgver  = "1.0.5"
)

func printUsage() {
	empty := "Empty. You need to set this up."
	tfToken := os.Getenv("TF_TOKEN")
	if tfToken == "" {
		tfToken = empty
	} else {
		tfToken = "__DELIBERATELY_REDACTED__"
	}
	tfDomain := os.Getenv("TF_DOMAIN")
	if tfDomain == "" {
		tfDomain = empty
	}
	tfOrg := os.Getenv("TF_ORG")
	if tfOrg == "" {
		tfOrg = empty
	}
	fmt.Printf(prgname + " v" + prgver + "\n" +
		"Terraform Enterprise utility - https://github.com/queone/tfe\n" +
		"Usage: " + prgname + " [options]\n" +
		"  -o [filter]              List orgs, filter option\n" +
		"  -m [filter]              List only latest version of modules, filter option\n" +
		"  -ma [filter]             List all version of modules, filter option\n" +
		"  -w [filter]              List workspaces (100 limit), filter option\n" +
		"  -ws WS_NAME              Show workspace details\n" +
		"  -wc WS_SRC WS_DES        Clone workspace named WS_SRC as WS_DES\n" +
		"  -?, -h, --help           Print this usage page\n" +
		"\n" +
		"  Note: This utility relies on below 3 critical environment variables:\n" +
		"    TF_TOKEN     A security token to access the respective TFE instance\n" +
		"    TF_DOMAIN    The TFE domain name (https://tfe.mydomain.com, etc)\n" +
		"    TF_ORG       The TFE Organization name (MY_ORG, etc)\n" +
		"\n" +
		"  Current values:\n" +
		"    TF_TOKEN=\"" + tfToken + "\"\n" +
		"    TF_DOMAIN=\"" + tfDomain + "\"\n" +
		"    TF_ORG=\"" + tfOrg + "\"\n")
	os.Exit(0)
}

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

func main() {
	numberOfArguments := len(os.Args[1:]) // Not including the program itself
	if numberOfArguments < 1 || numberOfArguments > 3 {
		printUsage() // Don't accept less than 1 or more than 3 arguments
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
		case "-ws":
			ShowWorkspace(client, os.Getenv("TF_ORG"), filter)
		}
	case 3: // Process 2-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		arg3 := os.Args[3]
		client := SetupClient()
		switch arg1 {
		case "-wc":
			CloneWorkspace(client, os.Getenv("TF_ORG"), arg2, arg3)
		}
	default:
		printUsage()
	}
}
