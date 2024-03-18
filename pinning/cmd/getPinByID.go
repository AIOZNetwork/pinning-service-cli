package cmd

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pinningcli "github.com/tue.phan/pinning-service-cli"
)

var (
	getPinByIDCmd = &cobra.Command{
		Use:   "get-pin",
		Short: "Retrieve a pin by its ID.",
		Long:  "Retrieve a pin by its ID.",
		Run: func(cmd *cobra.Command, args []string) {

			prepare()

			log.Debug("Running get-pin")

			var pinning pinningcli.PinningCli

			criteria := pinningcli.GetPinByIDCriteria{
				ID: id,
			}
			pin, errs := pinning.GetPinByID(ids, criteria)

			for _, err := range errs {
				log.Error(err)
			}

			if len(errs) > 0 {
				return
			}

			byteData, err := json.MarshalIndent(pin, "", "  ")
			if err != nil {
				log.Error(err)
				return
			}
			fmt.Println(string(byteData))
		},
	}
)

func init() {
	getPinByIDCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key")
	getPinByIDCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "Secret API Key")
	getPinByIDCmd.Flags().StringVarP(&id, "id", "", "", "Retrieves a pinned item from the IPFS network by ID.")
	rootCmd.AddCommand(getPinByIDCmd)
}
