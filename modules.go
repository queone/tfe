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

func ListModules(client *tfe.Client, orgName string, filter string, ver string) {
	// Lists registered modules, with a name filter option, which version option
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
				updatedAt, err := time.Parse(time.RFC3339, m.UpdatedAt)
				if err != nil {
					log.Fatalf("Error parsing updated timestamp: %v", err)
				}
				if utl.SubString(strings.ToLower(m.Name), filter) {
					for _, v := range m.VersionStatuses {
						fmt.Printf("%-80s %-10s %-06s %s\n", "localterraform.com/"+m.Namespace+"/"+m.Name+"/"+m.Provider, v.Version, v.Status, updatedAt.Format("2006-01-02 15:04"))
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
				updatedAt, err := time.Parse(time.RFC3339, m.UpdatedAt)
				if err != nil {
					log.Fatalf("Error parsing updated timestamp: %v", err)
				}
				if i, exists := latestVersions[m.Name]; exists {
					v := m.VersionStatuses[i]
					fmt.Printf("%-80s %-10s %-06s %s\n", "localterraform.com/"+m.Namespace+"/"+m.Name+"/"+m.Provider, v.Version, v.Status, updatedAt.Format("2006-01-02 15:04"))
				}
			}
		}
	}
}

func PrintModuleDetails(client *tfe.Client, modID tfe.RegistryModuleID) {
	// Prints details of the module
	module, err := client.RegistryModules.Read(context.Background(), modID)
	if err != nil {
		log.Fatalf("Error retrieving module with ID %v: %v", modID, err)
	}
	fmt.Printf("%+v\n", *module)
}
