package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/44smkn/s3select/pkg/build"
	"github.com/44smkn/s3select/pkg/cli"
	"github.com/44smkn/s3select/pkg/cli/root"
	"github.com/44smkn/s3select/pkg/log"
)

const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeInitializeError = 10 + iota
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
		return ExitCodeInitializeError
	}
	cliFactory, err := cli.NewFactory(buildVersion, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize process: %s", err.Error())
		return ExitCodeInitializeError
	}
	rootCmd := root.NewCmdRoot(cliFactory, buildVersion, buildDate)

	if _, err := rootCmd.ExecuteC(); err != nil {
		switch {
		case errors.Is(err, cli.ValidateConfigError):
			// advice for configuration
		}
		fmt.Fprintln(os.Stderr, err.Error())
	}

	return ExitCodeOK
}
