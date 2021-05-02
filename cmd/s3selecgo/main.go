package main

import (
	"fmt"
	"os"

	"github.com/44smkn/s3selecgo/pkg/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
	ExitCodeParseFlagsError
	ExitCodeLoggerError
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	cfg, err := loadConfig(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "コマンド引数のパースに失敗しました\n%v\n", err)
		return ExitCodeParseFlagsError
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ロガーの生成に失敗しました\n%v\n", err)
		return ExitCodeLoggerError
	}

	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion("ap-northeast-1")))
	svc := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(cfg.BucketName),
		Prefix: aws.String(cfg.KeyPrefix),
	}
	resp, err := svc.ListObjects(params)
	if err != nil {
		logger.Fatal("object listing is failed", zap.String("error", err.Error()))
	}
	for _, item := range resp.Contents {
		params := &s3.SelectObjectContentInput{
			Bucket:          aws.String(cfg.BucketName),
			Key:             item.Key,
			ExpressionType:  aws.String(s3.ExpressionTypeSql),
			Expression:      aws.String(cfg.SQL),
			RequestProgress: &s3.RequestProgress{},
			InputSerialization: &s3.InputSerialization{
				CompressionType: aws.String("GZIP"),
				CSV: &s3.CSVInput{
					FileHeaderInfo: aws.String(s3.FileHeaderInfoNone),
					FieldDelimiter: aws.String(" "),
				},
			},
			OutputSerialization: &s3.OutputSerialization{
				CSV: &s3.CSVOutput{
					FieldDelimiter: aws.String(" "),
				},
			},
		}
		resp, err := svc.SelectObjectContent(params)
		if err != nil {
			logger.Error("s3 select is failed", zap.String("error", err.Error()))
		}
		defer resp.EventStream.Close()

		for event := range resp.EventStream.Events() {
			// メッセージタイプ（イベントのタイプ）が ``Records`` の場合にメッセージからデータを取り出す
			v, ok := event.(*s3.RecordsEvent)
			if ok {
				fmt.Println(string(v.Payload))
			}
		}
	}
	return ExitCodeOK
}

func loadConfig(args []string) (config.FlagConfig, error) {
	cfg := config.FlagConfig{}
	fs := pflag.NewFlagSet("command line args", pflag.ExitOnError)
	cfg.BindFlags(fs)

	if err := fs.Parse(os.Args); err != nil {
		return config.FlagConfig{}, err
	}
	if err := cfg.Validate(); err != nil {
		return config.FlagConfig{}, err
	}

	return cfg, nil
}
