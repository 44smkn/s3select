package cliutil

import (
	"io"
	"os"

	"github.com/44smkn/s3select/pkg/aws"
	"github.com/44smkn/s3select/pkg/config"
	"github.com/mattn/go-colorable"
)

type Factory struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer

	AwsClient func() (aws.Cloud, error)
	Config    config.CliConfig

	Executable string
}

func NewFactory(appVersion string) (*Factory, error) {
	var cachedCfg *config.CliConfig
	var configErr error
	cfgFunc := func() (*config.CliConfig, error) {
		if cachedCfg != nil || configErr != nil {
			return cachedCfg, configErr
		}
		// parse
		return cachedCfg, configErr
	}

	executable := "s3select"
	if exe, err := os.Executable(); err == nil {
		executable = exe
	}

	return &Factory{
		In:     os.Stdin,
		Out:    colorable.NewColorable(os.Stdout),
		ErrOut: colorable.NewColorable(os.Stderr),

		AwsClient: func() (aws.Cloud, error) {
			cfg, err := cfgFunc()
			if err != nil {
				return nil, err
			}
			return aws.NewCloud(cfg.AWSConfig)
		},
		Executable: executable,
	}, nil
}
