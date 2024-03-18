package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pinningcli "github.com/tue.phan/pinning-service-cli"
)

const (
	credentialsDir  = "~/.pinning"
	credentialsFile = "credentials"
	apiKeyLength    = 24
	secretKeyLength = 44
)

func validateCredentials(apiKey string, secretKey string) error {
	if len(apiKey) != apiKeyLength || len(secretKey) != secretKeyLength {
		return fmt.Errorf("Invalid API key or secret key")
	}

	var pinning pinningcli.PinningCli

	result, errs := pinning.TestAuthentication(ids)
	for _, err := range errs {
		log.Error(err)
	}

	if len(errs) > 0 {
		return nil
	}
	if !result.Success {
		return fmt.Errorf("API key not valid")
	}

	return nil
}

var (
	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Store and access API key and secret key credentials securely.",
		Long:  `Store and access API key and secret key credentials securely in a credentials file.`,
		Run: func(cmd *cobra.Command, args []string) {

			prepare()

			if apiKey != "" && secretKey != "" {
				if err := validateCredentials(apiKey, secretKey); err != nil {
					log.Error(err)
					return
				}

				if err := saveCredentials(apiKey, secretKey); err != nil {
					if err.Error() == pinningcli.ErrNotOverwritten {
						log.Warn(pinningcli.ErrNotOverwritten)
						return
					}
					log.Error("Failed to save credentials: ", err)
					return
				}
				fmt.Println("Credentials saved successfully.")
			} else {
				apiKey, secretKey, err := readCredentials()
				if err != nil {
					fmt.Println("Failed to read credentials:", err)
					return
				}
				fmt.Println("API key:", apiKey)
				fmt.Println("Secret key:", secretKey)
			}

		},
	}
)

func init() {
	loginCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key (MANDATORY)")
	loginCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "Secret API Key (MANDATORY)")
	rootCmd.AddCommand(loginCmd)
}

func saveCredentials(apiKey, secretKey string) error {
	credentialsDir, err := expandHomeDir(credentialsDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(credentialsDir, 0700)
	if err != nil {
		return err
	}

	credentialsFile := filepath.Join(credentialsDir, credentialsFile)

	if _, err := os.Stat(credentialsFile); err == nil {
		// Credentials file already exists
		fmt.Print("Credentials already exist. Do you want to overwrite? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			return errors.New(pinningcli.ErrNotOverwritten)
		}
	}

	credentials := map[string]string{
		"APIKey":    apiKey,
		"SecretKey": secretKey,
	}

	data, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(credentialsFile, data, 0600)
	if err != nil {
		return err
	}

	return nil
}

func readCredentials() (string, string, error) {
	credentialsDir, err := expandHomeDir(credentialsDir)
	if err != nil {
		return "", "", err
	}

	credentialsFile := filepath.Join(credentialsDir, credentialsFile)

	data, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return "", "", err
	}

	var credentials map[string]string
	err = json.Unmarshal(data, &credentials)
	if err != nil {
		return "", "", err
	}

	apiKey := credentials["APIKey"]
	secretKey := credentials["SecretKey"]

	return apiKey, secretKey, nil
}

func expandHomeDir(path string) (string, error) {
	if path[:2] == "~/" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(homeDir, path[2:])
	}
	return path, nil
}
