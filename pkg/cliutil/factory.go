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
	Config    config.Config

	Executable string
}

func NewFactory(appVersion string) (*Factory, error) {
	var cachedCfg config.Config
	var cfgErr error
	cfgFunc := func() (config.Config, error) {
		if cachedCfg != nil || cfgErr != nil {
			return cachedCfg, cfgErr
		}
		cachedCfg, cfgErr = config.ParseDefaultConfig()
		return cachedCfg, cfgErr
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
			cloudCfg := aws.CloudConfig{
				Region: cfg.GetAWSRegion(),
			}
			return aws.NewCloud(cloudCfg)
		},
		Executable: executable,
	}, nil
}
