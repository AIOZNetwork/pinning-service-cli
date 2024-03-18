package cmd

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pinningcli "github.com/tue.phan/pinning-service-cli"
)

var (
	listPinsCmd = &cobra.Command{
		Use:   "list-pins",
		Short: "List of user pins.",
		Long:  `This endpoint returns data on what content the sender has pinned to IPFS through Pinning.`,
		Run: func(cmd *cobra.Command, args []string) {

			prepare()

			log.Debug("Running list-pins")

			var pinning pinningcli.PinningCli

			keyvaluesquery, errs := parseKeyValues(keyvaluesquery)
			for _, err := range errs {
				log.Error(err)
			}

			if len(errs) > 0 {
				return
			}

			criteria := pinningcli.ListPinsCriteria{
				Offset:         offset,
				Limit:          limit,
				SortBy:         sortBy,
				SortOrder:      sortOrder,
				Pinned:         pinned,
				KeyValuesQuery: keyvaluesquery,
			}
			listPinResult, errs := pinning.GetPinsList(ids, criteria)

			for _, err := range errs {
				log.Error(err)
			}

			if len(errs) > 0 {
				return
			}

			byteData, err := json.MarshalIndent(listPinResult, "", "  ")
			if err != nil {
				log.Error(err)
				return
			}
			fmt.Println(string(byteData))

		},
	}
)

func init() {
	listPinsCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key")
	listPinsCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "Secret API Key")
	listPinsCmd.Flags().UintVarP(&offset, "offset", "", 0, "(default 0)")
	listPinsCmd.Flags().UintVarP(&limit, "limit", "", 10, "(default 10)")
	listPinsCmd.Flags().StringVarP(&pinned, "pinned", "", "", "Filter by pinned status (options: all, true, false) (default all)")
	listPinsCmd.Flags().StringVarP(&sortBy, "sortBy", "", "", "Field to sort by (options: created_at, size, name). Defaults to created_at.")
	listPinsCmd.Flags().StringVarP(&sortOrder, "sortOrder", "", "", "Sort direction (options: ASC, DESC). Defaults to DESC.")
	listPinsCmd.Flags().StringVarP(&name, "name", "", "", "Name of the hash to search")
	listPinsCmd.Flags().StringArrayVarP(&keyvaluesquery, "keyvalue", "", nil, "Additional key/value query")
	rootCmd.AddCommand(listPinsCmd)
}
