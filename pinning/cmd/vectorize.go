package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pinningcli "github.com/tue.phan/pinning-service-cli"
)

var (
	vectorizeCmd = &cobra.Command{
		Use:   "vectorize",
		Short: "Vectorize pinned files and query semantic vectors",
		Long:  "Vectorize pinned files and query semantic vectors.",
	}

	vectorizeFileCmd = &cobra.Command{
		Use:   "file",
		Short: "Vectorize a pinned file",
		Long:  "Vectorize a pinned file.",
		Run: func(cmd *cobra.Command, args []string) {
			prepare()

			if cid == "" {
				log.Error("cid is mandatory")
				os.Exit(1)
			}

			var pinning pinningcli.PinningCli
			result, errs := pinning.VectorizeFile(ids, pinningcli.VectorizeFileCriteria{CID: cid})
			printResult(result, errs)
		},
	}

	vectorizeDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Soft delete a vectorized file",
		Long:  "Soft delete a vectorized file.",
		Run: func(cmd *cobra.Command, args []string) {
			prepare()

			if cid == "" {
				log.Error("cid is mandatory")
				os.Exit(1)
			}

			var pinning pinningcli.PinningCli
			result, errs := pinning.DeleteVectorizedFile(ids, pinningcli.DeleteVectorizedFileCriteria{CID: cid})
			printResult(result, errs)
		},
	}

	vectorizeQueryCmd = &cobra.Command{
		Use:   "query",
		Short: "Query semantic vectors",
		Long:  "Query semantic vectors.",
		Run: func(cmd *cobra.Command, args []string) {
			prepare()

			if queryText == "" {
				log.Error("query text is mandatory")
				os.Exit(1)
			}

			var pinning pinningcli.PinningCli
			result, errs := pinning.QueryVectors(ids, pinningcli.QueryVectorsCriteria{
				Text: queryText,
				TopK: topK,
			})
			printResult(result, errs)
		},
	}
)

func printResult(result interface{}, errs []error) {
	for _, err := range errs {
		log.Error(err)
	}

	if len(errs) > 0 {
		os.Exit(1)
	}

	byteData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	fmt.Println(string(byteData))
}

func init() {
	vectorizeFileCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key")
	vectorizeFileCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "Secret API Key")
	vectorizeFileCmd.Flags().StringVarP(&cid, "cid", "", "", "File CID to vectorize")

	vectorizeDeleteCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key")
	vectorizeDeleteCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "Secret API Key")
	vectorizeDeleteCmd.Flags().StringVarP(&cid, "cid", "", "", "File CID to soft delete from vector search")

	vectorizeQueryCmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key")
	vectorizeQueryCmd.Flags().StringVarP(&secretKey, "secret", "s", "", "Secret API Key")
	vectorizeQueryCmd.Flags().StringVarP(&queryText, "text", "t", "", "Semantic query text")
	vectorizeQueryCmd.Flags().UintVarP(&topK, "top-k", "", 10, "Maximum number of vector search results (max 100)")

	vectorizeCmd.AddCommand(vectorizeFileCmd)
	vectorizeCmd.AddCommand(vectorizeDeleteCmd)
	vectorizeCmd.AddCommand(vectorizeQueryCmd)
	rootCmd.AddCommand(vectorizeCmd)
}
