package main

import (
	"fmt"
	"os"

	"github.com/queone/utl"
)

const (
	program_name    = "tfe"
	program_version = "1.2.0"
	config_file     = "config.yaml"
)

func printUsage() {
	n := utl.Whi2(program_name)
	v := program_version
	cfgFile := utl.Yel("~/.config/" + program_name + "/config.yaml")
	usage := fmt.Sprintf("%s v%s\n"+
		"Terraform Cloud CLI utility — https://github.com/queone/tfe\n"+
		"%s\n"+
		"  %s [options] [arguments]\n"+
		"\n"+
		"%s\n"+
		"  This utility provides basic functionality for interacting with a "+
		"Terraform Enterprise or TF Cloud account instance.\n"+
		"\n"+
		"%s\n"+
		"  You can authenticate using one of two methods:\n"+
		"  1. Set the following environment variables:\n"+
		"     • TF_ORG=<value>     TFE Organization name (e.g. MYORG)\n"+
		"     • TF_DOMAIN=<value>  TFE domain name (e.g. https://app.terraform.io)\n"+
		"     • TF_TOKEN=<value>   Security token to access the TFE instance\n"+
		"  2. Create a YAML configuration file at %s with those variables:\n"+
		"     • TF_ORG:    <value>\n"+
		"     • TF_DOMAIN: <value>\n"+
		"     • TF_TOKEN:  <value>\n"+
		"\n"+
		"%s\n"+
		"  The following options are available:\n"+
		"  -o [filter]       List organizations; filter option\n"+
		"  -m[j] [filter]    List only latest version of modules; filter option; JSON option\n"+
		"  -ma [filter]      List all versions of modules; filter option\n"+
		"  -w [filter]       List workspaces (100 limit); filter option\n"+
		"  -ws NAME          Show workspace details\n"+
		"  -wc SRC DEST      Clone workspace named SRC as a new workspace named DEST\n"+
		"  -?, -h, --help    Show this help message and exit\n"+
		"\n"+
		"%s\n"+
		"  %s -o myorg\n"+
		"  %s -m mymodule\n"+
		"  %s -ws myworkspace\n"+
		"  %s -wc source_ws_name new_ws_name\n"+
		"  %s -h\n",
		n, v, utl.Whi2("Usage"), n, utl.Whi2("Overview"), utl.Whi2("Authentication"), cfgFile, utl.Whi2("Options"), utl.Whi2("Examples"), n, n, n, n, n)
	fmt.Print(usage)
	os.Exit(0)
}

func main() {
	numberOfArguments := len(os.Args[1:]) // Don't include the program itself
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
		client, tfOrg := setupClient()
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
		client, tfOrg := setupClient()
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
		client, tfOrg := setupClient()
		switch arg1 {
		case "-wc":
			CloneWorkspace(client, tfOrg, arg2, arg3)
		}
	default:
		printUsage()
	}
}
