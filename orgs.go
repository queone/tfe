package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-tfe"
	"github.com/queone/utl"
)

func ListOrganizations(client *tfe.Client, filter string) {
	// Lists organizations, with a name filter option
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
