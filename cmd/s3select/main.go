package main

import (
	"fmt"
	"os"

	"github.com/44smkn/s3select/pkg/build"
	"github.com/44smkn/s3select/pkg/cli/root"
	"github.com/44smkn/s3select/pkg/cliutil"
	"github.com/44smkn/s3select/pkg/log"
)

const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
	ExitCodeParseFlagsError
	ExitCodeLoggerError
	ExitCodeCloudError
	ExitCodeObjectListingError
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	buildDate := build.Date
	buildVersion := build.Version

	logLevel := "info"
	if os.Getenv("DEBUG") == "true" {
		logLevel = "debug"
	}
	logger, err := log.New(logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate logger: %s", err.Error())
	}
	cliFactory, err := cliutil.NewFactory(buildVersion, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize process: %s", err.Error())
	}
	rootCmd := root.NewCmdRoot(cliFactory, buildVersion, buildDate)

	if cmd, err := rootCmd.ExecuteC(); err != nil {
		// TODO: error handling
		fmt.Fprintln(os.Stderr, cmd.UsageString())
	}

	return ExitCodeOK
}
