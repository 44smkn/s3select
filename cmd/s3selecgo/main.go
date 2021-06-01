package main

import (
	"context"
	"fmt"
	"os"

	"github.com/44smkn/s3selecgo/pkg/aws"
	"github.com/44smkn/s3selecgo/pkg/config"
	"github.com/44smkn/s3selecgo/pkg/log"
	"github.com/44smkn/s3selecgo/pkg/query"
	awssdk "github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/pflag"
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
	cfg, err := loadConfig(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse command args:\n%v\n", err)
		return ExitCodeParseFlagsError
	}
	logger, err := log.NewLogger(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger instance: \n%v\n", err)
		return ExitCodeLoggerError
	}

	cloud, err := aws.NewCloud(cfg.AWSConfig)
	if err != nil {
		return ExitCodeCloudError
	}

	ctx := context.Background()

	req := s3sdk.ListObjectsV2Input{
		Bucket: awssdk.String(cfg.S3SelectConfig.BucketName),
		Prefix: awssdk.String(cfg.S3SelectConfig.KeyPrefix),
	}
	objects, err := cloud.S3().ListObjectsV2AsList(ctx, &req)
	if err != nil {
		logger.Sugar().Errorf("failed to execute s3listobeject api: %s", err.Error())
		return ExitCodeObjectListingError
	}

	sq := query.NewDefaultStorageQueryer(cfg.S3SelectConfig, cloud)
	for _, o := range objects {
		sq.Query(ctx, cfg.S3SelectConfig.BucketName, *o.Key, os.Stdout)
	}
	return ExitCodeOK
}

func loadConfig(args []string) (config.CliConfig, error) {
	cfg := config.CliConfig{}
	fs := pflag.NewFlagSet("command line args", pflag.ExitOnError)
	cfg.BindFlags(fs)

	if err := fs.Parse(os.Args); err != nil {
		return config.CliConfig{}, err
	}
	if err := cfg.S3SelectConfig.Validate(); err != nil {
		return config.CliConfig{}, err
	}

	return cfg, nil
}
