package config

import (
	"github.com/44smkn/s3selecgo/pkg/aws"
	"github.com/spf13/pflag"
)

const (
	flagLogLevel    = "log-level"
	defaultLogLevel = "info"
)

type CliConfig struct {
	LogLevel  string
	S3Config  S3Config
	AWSConfig aws.CloudConfig
}

// BindFlags binds the command line flags to the fields in the config object
func (cfg *CliConfig) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cfg.LogLevel, flagLogLevel, defaultLogLevel, "Set the cli log level - info(default), debug")
	cfg.S3Config.BindFlags(fs)
	cfg.AWSConfig.BindFlags(fs)
}
