package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-tfe"
)

const (
	prgname = "tfe"
	prgver  = "1.1.2"
)

func printUsage() {
	empty := "Empty. You need to set this up."
	tfOrg := os.Getenv("TF_ORG")
	tfDomain := os.Getenv("TF_DOMAIN")
	tfToken := os.Getenv("TF_TOKEN")
	if tfOrg == "" {
		tfOrg = empty
	}
	if tfDomain == "" {
		tfDomain = empty
	}
	if tfToken == "" {
		tfToken = empty
	} else {
		tfToken = "__DELIBERATELY_REDACTED__"
	}
	fmt.Printf(prgname + " v" + prgver + "\n" +
		"Terraform cloud utility - https://github.com/queone/tfe\n" +
		"Usage: " + prgname + " [options]\n" +
		"  -o [filter]              List orgs, filter option\n" +
		"  -m[j] [filter]           List only latest version of modules, filter option; JSON option\n" +
		"  -ma [filter]             List all version of modules, filter option\n" +
		"  -w [filter]              List workspaces (100 limit), filter option\n" +
		"  -ws NAME                 Show workspace details\n" +
		"  -wc SRC DES              Clone workspace named SRC as DES\n" +
		"  -?, -h, --help           Print this usage page\n" +
		"\n" +
		"  Note: This utility relies on below 3 critical environment variables:\n" +
		"    TF_ORG       TFE Organization name (MYORG, etc)\n" +
		"    TF_DOMAIN    TFE domain name (https://app.terraform.io, etc)\n" +
		"    TF_TOKEN     Security token to access the respective TFE instance\n" +
		"\n" +
		"  Current values:\n" +
		"    TF_ORG=\"" + tfOrg + "\"\n" +
		"    TF_DOMAIN=\"" + tfDomain + "\"\n" +
		"    TF_TOKEN=\"" + tfToken + "\"\n")
	os.Exit(0)
}

func SetupClient(tfOrg, tfDomain, tfToken string) *tfe.Client {
	// Check if essential environment variables are valid
	if tfToken == "" || tfOrg == "" || tfDomain == "" {
		log.Fatal("One or more required environment variables (TF_TOKEN, TF_ORG, TF_DOMAIN) are not set.")
	}

	// Set up a configuration with the API token
	config := &tfe.Config{
		Token:   tfToken,
		Address: tfDomain,
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

	// Retrieve the 3 essential environment variables
	tfOrg := os.Getenv("TF_ORG")
	tfDomain := os.Getenv("TF_DOMAIN")
	tfToken := os.Getenv("TF_TOKEN")

	switch numberOfArguments {
	case 1: // Process 1-argument requests
		arg1 := os.Args[1]
		switch arg1 {
		case "-?", "-h", "--help":
			printUsage()
		}
		client := SetupClient(tfOrg, tfDomain, tfToken)
		switch arg1 {
		case "-o":
			ListOrganizations(client, "")
		case "-m":
			ListModules(client, tfOrg, "", "latest")
		case "-ma":
			ListModules(client, tfOrg, "", "all")
		case "-mj":
			ListModules(client, tfOrg, "", "json")
		case "-w":
			ListWorkspaces(client, tfOrg, "")
		}
	case 2: // Process 2-argument requests
		arg1 := os.Args[1]
		filter := os.Args[2]
		client := SetupClient(tfOrg, tfDomain, tfToken)
		switch arg1 {
		case "-o":
			ListOrganizations(client, filter)
		case "-m":
			ListModules(client, tfOrg, filter, "latest")
		case "-ma":
			ListModules(client, tfOrg, filter, "all")
		case "-mj":
			ListModules(client, tfOrg, filter, "json")
		case "-w":
			ListWorkspaces(client, tfOrg, filter)
		case "-ws":
			ShowWorkspace(client, tfOrg, filter)
		}
	case 3: // Process 2-argument requests
		arg1 := os.Args[1]
		arg2 := os.Args[2]
		arg3 := os.Args[3]
		client := SetupClient(tfOrg, tfDomain, tfToken)
		switch arg1 {
		case "-wc":
			CloneWorkspace(client, tfOrg, arg2, arg3)
		}
	default:
		printUsage()
	}
}
