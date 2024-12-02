package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-tfe"
	"github.com/queone/utl"
)

func GetRegistryModuleID(tfOrg string, mod *tfe.RegistryModule) tfe.RegistryModuleID {
	// Returns the RegistryModuleID type struct for this module
	if mod.RegistryName == "private" {
		return tfe.NewPrivateRegistryModuleID(tfOrg, mod.Name, mod.Provider)
	} else {
		return tfe.NewPublicRegistryModuleID(tfOrg, mod.Namespace, mod.Name, mod.Provider)
	}
}

func GetModuleVersion(client *tfe.Client, tfOrg string, mod *tfe.RegistryModule, verStr string) *tfe.RegistryModuleVersion {
	// Returns RegistryModuleVersion struct for this module's verStr version
	moduleID := GetRegistryModuleID(tfOrg, mod)
	ver, err := client.RegistryModules.ReadVersion(context.Background(), moduleID, verStr)
	if err != nil {
		log.Fatalf("Error fetching module %s version %s: %v", mod.Name, verStr, err)
	}
	return ver
}

func ListModules(client *tfe.Client, tfOrg string, filter string, qualifier string) {
	// List registered modules, with a name filter option, and qualifier option
	// to print all version, latest version or in JSON format if it's a single module.
	options := tfe.RegistryModuleListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
	}
	var allModules []*tfe.RegistryModule
	var matchingModules []*tfe.RegistryModule

	// Retrieve all modules from the organization
	for {
		modules, err := client.RegistryModules.List(context.Background(), tfOrg, &options)
		if err != nil {
			log.Fatalf("Error listing modules for organization %s: %v", tfOrg, err)
		}
		allModules = append(allModules, modules.Items...)
		if modules.NextPage == 0 {
			break
		}
		options.PageNumber = modules.NextPage
	}

	// Create list of matching modules
	for _, mod := range allModules {
		name := strings.ToLower(mod.Name)
		if utl.SubString(name, filter) {
			matchingModules = append(matchingModules, mod)
		}
	}

	if len(matchingModules) == 1 {
		if qualifier == "json" {
			PrintSingleModuleDetails(client, tfOrg, matchingModules[0], true)
		} else {
			PrintSingleModuleDetails(client, tfOrg, matchingModules[0], false)
		}
	} else if len(matchingModules) > 1 {
		if qualifier == "all" {
			// Print each version of all module names matching filter
			for _, mod := range matchingModules {
				modNamespace := "localterraform.com/" + mod.Namespace + "/" + mod.Name + "/" + mod.Provider
				for _, v := range mod.VersionStatuses {
					ver := GetModuleVersion(client, tfOrg, mod, v.Version)
					updatedAt, _ := time.Parse(time.RFC3339, ver.UpdatedAt)
					updated_at := updatedAt.Format("2006-Jan-02 15:04")
					fmt.Printf("%-80s %-10s %s\n", modNamespace, v.Version, updated_at)
				}
			}
		} else {
			// Map to store the latest version of each module
			latestVersions := make(map[string]int)

			// Iterate over all modules and track the latest version
			for _, mod := range matchingModules {
				for i, v := range mod.VersionStatuses {
					if current, exists := latestVersions[mod.Name]; !exists ||
						strings.Compare(v.Version, mod.VersionStatuses[current].Version) > 0 {
						latestVersions[mod.Name] = i
					}
				}
			}

			// Print the latest version of each filtered module
			for _, mod := range matchingModules {
				modNamespace := "localterraform.com/" + mod.Namespace + "/" + mod.Name + "/" + mod.Provider
				if i, exists := latestVersions[mod.Name]; exists {
					v := mod.VersionStatuses[i]
					ver := GetModuleVersion(client, tfOrg, mod, v.Version)
					updatedAt, _ := time.Parse(time.RFC3339, ver.UpdatedAt)
					updated_at := updatedAt.Format("2006-Jan-02 15:04")
					fmt.Printf("%-80s %-10s %s\n", modNamespace, v.Version, updated_at)
				}
			}
		}
	}
}

func PrintSingleModuleDetails(client *tfe.Client, tfOrg string, mod *tfe.RegistryModule, json bool) {
	// Print all details related to this one specific module.
	// The tfe.Client arg will be used to print other details in the future.
	if json {
		// Print the full JSON object
		utl.PrintJsonColor(mod)
	} else {
		// Prints most important attributes in YAML-like format
		fmt.Printf("%s\n", utl.Gra("# Terraform Cloud Registry Module"))
		fmt.Printf("%s: %s\n", utl.Blu("id"), utl.Gre(mod.ID))
		fmt.Printf("%s: %s\n", utl.Blu("name"), utl.Gre(mod.Name))
		fmt.Printf("%s: %s\n", utl.Blu("provider"), utl.Gre(mod.Provider))
		fmt.Printf("%s: %s\n", utl.Blu("namespace"), utl.Gre(mod.Namespace))

		// Capture/reformat created/updated at dates from "2024-12-01T17:00:58.518Z"
		createdAt, _ := time.Parse(time.RFC3339, mod.CreatedAt)
		created_at := createdAt.Format("2006-Jan-02 15:04")
		updatedAt, _ := time.Parse(time.RFC3339, mod.UpdatedAt)
		updated_at := updatedAt.Format("2006-Jan-02 15:04")
		fmt.Printf("%s: %s\n", utl.Blu("created_at"), utl.Gre(created_at))
		fmt.Printf("%s: %s\n", utl.Blu("updated_at"), utl.Gre(updated_at))

		if mod.VCSRepo != nil {
			fmt.Printf("%s: %s\n", utl.Blu("repo_url"), utl.Gre(mod.VCSRepo.RepositoryHTTPURL))
		}
		if len(mod.VersionStatuses) > 0 {
			fmt.Printf("%s:\n", utl.Blu("versions"))
			for _, v := range mod.VersionStatuses {
				ver := GetModuleVersion(client, tfOrg, mod, v.Version)
				//utl.PrintJsonColor(ver) // DEBUG
				createdAt, _ := time.Parse(time.RFC3339, ver.CreatedAt)
				created_As := createdAt.Format("2006-Jan-02 15:04")
				updatedAt, _ := time.Parse(time.RFC3339, ver.UpdatedAt)
				updated_at := updatedAt.Format("2006-Jan-02 15:04")
				fmt.Printf("  %-26s %-8s %-20s %s\n", ver.ID, v.Version, created_As, updated_at)
			}

		}
	}
}
