package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/tue.phan/pinning-service-cli/pinning/cmd"
)

func main() {
	log.SetOutput(os.Stderr)

	if err := cmd.Execute(); err != nil {
		log.Fatal("Fatal error during execution")
		os.Exit(1)
	}

	log.Debug("Command completed")
}
