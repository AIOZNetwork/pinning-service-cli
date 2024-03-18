package cmd

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pinningcli "github.com/tue.phan/pinning-service-cli"
)

var (
	unpinCmd = &cobra.Command{
		Use:   "unpin",
		Short: "Unpin a hash from Pinning",
		Long:  `Unpin a hash previously pinned to Pinning.`,
		Run: func(cmd *cobra.Command, args []string) {

			prepare()

			log.Debug("Running unpin")

			if id == "" {
				log.Error("id to unpin is mandatory")
				return
			}

			var pinning pinningcli.PinningCli

			criteria := pinningcli.UnpinCriteria{
				ID: id,
			}
			result, errs := pinning.Unpin(ids, criteria)

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
	unpinCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key (MANDATORY)")
	unpinCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "Secret API Key (MANDATORY)")
	unpinCmd.Flags().StringVarP(&id, "id", "", "", "id to unpin IPFS (MANDATORY)")
	rootCmd.AddCommand(unpinCmd)
}
