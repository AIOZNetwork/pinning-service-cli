package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pinningcli "github.com/tue.phan/pinning-service-cli"
	"gorm.io/datatypes"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

var (
	apiKey         string
	secretKey      string
	hash           string
	file           string
	name           string
	offset         uint
	limit          uint
	keyvaluesquery []string
	sortBy         string
	sortOrder      string
	keyvalues      []string
	pinned         string
	id             string

	ids pinningcli.PinningIdentifiers

	rootCmd = &cobra.Command{
		Use:     "pinning",
		Short:   "pinning handles operations on pinning",
		Long:    `pinning handles operations on pinning`,
		Version: "0.0.1",
	}
)

func prepare() {
	sourceCredentials()
}

func sourceCredentials() {
	var value string

	// source credentials from environment
	k, s, _ := readCredentials()
	ids.ApiKey = k
	ids.SecretKey = s

	// source credentials from flags
	value = apiKey
	if value != "" {
		ids.ApiKey = value
	}
	value = secretKey
	if value != "" {
		ids.SecretKey = value
	}

	// Check that credentials have been provided
	if ids.ApiKey == "" {
		log.Error("API Key is mandatory")
		return
	}

	if ids.SecretKey == "" {
		log.Error("Secret Key is mandatory")
		return
	}
}

func parseKeyValues(keyvalues []string) (datatypes.JSON, []error) {
	var metadataJs map[string]string
	if len(keyvalues) > 0 {
		metadataJs = make(map[string]string)
	}
	var errors []error
	for _, keyvalue := range keyvalues {
		tokens := strings.Split(keyvalue, ":")
		if len(tokens) != 2 {
			errors = append(errors, fmt.Errorf("Failed to split keyvalue '%s'. Expected format: <key>:<value>", keyvalue))
		} else if _, ok := metadataJs[tokens[0]]; ok {
			errors = append(errors, fmt.Errorf("Token '%s' is defined twice", tokens[0]))
		} else {
			metadataJs[tokens[0]] = tokens[1]
		}
	}

	jsonData, err := json.Marshal(metadataJs)
	if err != nil {
		return nil, errors
	}
	return datatypes.JSON(jsonData), errors
}

func Execute() error {
	return rootCmd.Execute()
}
