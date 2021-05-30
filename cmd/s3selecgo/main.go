package main

import (
	"context"
	"fmt"
	"os"

	"github.com/44smkn/s3selecgo/pkg/aws"
	"github.com/44smkn/s3selecgo/pkg/config"
	"github.com/44smkn/s3selecgo/pkg/log"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
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

	for _, o := range objects {
		params := &s3sdk.SelectObjectContentInput{
			Bucket:          awssdk.String(cfg.S3SelectConfig.BucketName),
			Key:             o.Key,
			ExpressionType:  awssdk.String(s3.ExpressionTypeSql),
			Expression:      awssdk.String(cfg.S3SelectConfig.SQLExpression),
			RequestProgress: &s3.RequestProgress{},
			InputSerialization: &s3.InputSerialization{
				CompressionType: awssdk.String(cfg.S3SelectConfig.InputCompressionType),
				CSV: &s3.CSVInput{
					FileHeaderInfo: awssdk.String(s3.FileHeaderInfoNone),
					FieldDelimiter: awssdk.String(" "),
				},
			},
			OutputSerialization: &s3.OutputSerialization{
				CSV: &s3.CSVOutput{
					FieldDelimiter: awssdk.String(" "),
				},
			},
		}
		resp, err := cloud.S3().SelectObjectContent(params)
		if err != nil {
			logger.Error("failed to execute s3select api", zap.String("error", err.Error()))
		}
		defer resp.EventStream.Close()

		for event := range resp.EventStream.Events() {
			// If the event type is `records`, it fetch the data from the message.
			v, ok := event.(*s3.RecordsEvent)
			if ok {
				fmt.Println(string(v.Payload))
			}
		}
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
