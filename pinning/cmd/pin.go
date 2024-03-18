package cmd

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pinningcli "github.com/tue.phan/pinning-service-cli"
)

var (
	pinCmd = &cobra.Command{
		Use:   "pin",
		Short: "Add a hash to Pinning for asynchronous pinning",
		Long:  `Add a hash to Pinning for asynchronous pinning. Content added through this endpoint is pinned in the background and will show up in your pinned items once the content has been found / pinned. For this operation to succeed, the content for the hash you provide must already be pinned by another node on the IFPS network.`,
		Run: func(cmd *cobra.Command, args []string) {

			prepare()

			if hash == "" && file == "" {
				log.Error("Either hash or File to pin is mandatory")
				return
			}

			var pinning pinningcli.PinningCli

			keyvalues, errs := parseKeyValues(keyvalues)

			for _, err := range errs {
				log.Error(err)
			}

			if len(errs) > 0 {
				return
			}

			var result interface{}
			if file != "" {
				criteria := pinningcli.PinFileCriteria{
					File:              file,
					Name:              name,
					KeyValues:         keyvalues,
				}
				result, errs = pinning.PinFile(ids, criteria)
			} else {
				criteria := pinningcli.PinByHashCriteria{
					HashToPin: hash,
					Name:      name,
					KeyValues: keyvalues,
				}
				result, errs = pinning.PinByHash(ids, criteria)
			}

			for _, err := range errs {
				log.Error(err)
			}

			if len(errs) > 0 {
				return
			}

			byteData, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				log.Error(err)
				return
			}
			fmt.Println(string(byteData))

		},
	}
)

func init() {
	pinCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key")
	pinCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "Secret API Key")
	pinCmd.Flags().StringVarP(&hash, "hash", "", "", "Hash to pin to IPFS (MANDATORY)")
	pinCmd.Flags().StringVarP(&file, "file", "", "", "File or directory to pin to IPFS")
	pinCmd.Flags().StringVarP(&name, "name", "", "", "A name for this hash, for display in Pinning only")
	pinCmd.Flags().StringArrayVarP(&keyvalues, "keyvalue", "", nil, "Additional key/value to store in Pinning")
	rootCmd.AddCommand(pinCmd)
}
