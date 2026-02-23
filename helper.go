package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-tfe"
	"github.com/queone/utl"
)

func validateCredentials(creds map[string]string) error {
	for key, value := range creds {
		if value == "" {
			return fmt.Errorf("%s is empty — cannot login", key)
		}
	}
	return nil
}

func getLoginCredentials() (tfOrg, tfDomain, tfToken string, err error) {
	const (
		envVarTFORG    = "TF_ORG"
		envVarTFDOMAIN = "TF_DOMAIN"
		envVarTFTOKEN  = "TF_TOKEN"
	)

	// Try reading environment variables
	tfOrg = os.Getenv(envVarTFORG)
	tfDomain = os.Getenv(envVarTFDOMAIN)
	tfToken = os.Getenv(envVarTFTOKEN)

	// If any environment variable is set, use them
	if tfOrg != "" || tfDomain != "" || tfToken != "" {
		// Validate environment variables
		creds := map[string]string{
			"TF_ORG":    tfOrg,
			"TF_DOMAIN": tfDomain,
			"TF_TOKEN":  tfToken,
		}
		if err := validateCredentials(creds); err != nil {
			return "", "", "", err
		}
		return tfOrg, tfDomain, tfToken, nil
	}

	// Get credentials from config file
	fmt.Printf("One or more of the 3 environment variables login credentials are not set.\n" +
		"Checking config file credentials...\n")

	// Set up configuration directory path, and create if it doesn't exist
	confDir := filepath.Join(os.Getenv("HOME"), ".config", program_name)
	if utl.FileNotExist(confDir) {
		if err := os.Mkdir(confDir, 0700); err != nil {
			return "", "", "", fmt.Errorf("config directory %s does not exist — error trying "+
				"to create it for future use: %w", utl.Yel(confDir), err)
		}
	}

	filePath := filepath.Join(confDir, config_file)
	contents := "TF_ORG:     # TFE Organization name (MYORG, etc)\n" +
		"TF_DOMAIN:  # TFE domain name (https://app.terraform.io, etc)\n" +
		"TF_TOKEN:   # Security token to access the respective TFE instance\n"
	if utl.FileNotExist(filePath) {
		fmt.Printf("Config file %s does not exist. Creating a new one...\n", utl.Yel(filePath))
		if err := os.WriteFile(filePath, []byte(contents), 0600); err != nil {
			return "", "", "", fmt.Errorf("error creating config file: %w", err)
		}
		fmt.Println("Config file created. Please fill in the credentials and try again.")
		return "", "", "", fmt.Errorf("config file created, but credentials are empty")
	} else if utl.FileSize(filePath) < 1 {
		fmt.Printf("Config file %s is empty. Initializing it...\n", utl.Yel(filePath))
		if err := os.WriteFile(filePath, []byte(contents), 0600); err != nil {
			return "", "", "", fmt.Errorf("error initializing config file: %w", err)
		}
		fmt.Println("Config file initialized. Please fill in the credentials and try again.")
		return "", "", "", fmt.Errorf("config file initialized, but credentials are empty")
	}

	credsRaw, err := utl.LoadFileYaml(filePath)
	if err != nil {
		return "", "", "", fmt.Errorf("error reading YAML config file: %w", err)
	}
	creds := credsRaw.(map[string]any)
	tfOrg = utl.Str(creds["TF_ORG"])
	tfDomain = utl.Str(creds["TF_DOMAIN"])
	tfToken = utl.Str(creds["TF_TOKEN"])
	if err := validateCredentials(map[string]string{
		"TF_ORG":    tfOrg,
		"TF_DOMAIN": tfDomain,
		"TF_TOKEN":  tfToken,
	}); err != nil {
		return "", "", "", err
	}

	return tfOrg, tfDomain, tfToken, nil
}

// Sets up the TFE client and the default TF organization login
func setupClient() (*tfe.Client, string) {
	tfOrg, tfDomain, tfToken, err := getLoginCredentials()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Use the API token to create a new TFE configuration and instance
	config := &tfe.Config{
		Token:   tfToken,
		Address: tfDomain,
	}
	client, err := tfe.NewClient(config)
	if err != nil {
		fmt.Printf("Error creating TFE client: %v\n", err)
		os.Exit(1)
	}

	return client, tfOrg
}
