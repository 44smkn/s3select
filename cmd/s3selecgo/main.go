package main

import (
	"fmt"
	"os"

	"github.com/44smkn/s3selecgo/pkg/build"
	"github.com/44smkn/s3selecgo/pkg/cli/root"
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

	rootCmd := root.NewCmdRoot(buildVersion, buildDate)

	if cmd, err := rootCmd.ExecuteC(); err != nil {
		// TODO: error handling
		fmt.Fprintln(os.Stderr, cmd.UsageString())
	}

	// TODO: read config file

	return ExitCodeOK
}
