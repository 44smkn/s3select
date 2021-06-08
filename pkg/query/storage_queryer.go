package query

import (
	"context"
	"io"
	"strings"

	"github.com/44smkn/s3select/pkg/aws"
	"github.com/44smkn/s3select/pkg/config"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type StorageQueryer interface {
	Query(ctx context.Context, bucketName, objectKey string, writer io.Writer)
}

func NewQueryer(cfg config.S3SelectConfig, cloud aws.Cloud) (StorageQueryer, error) {
	switch et := strings.ToLower(cfg.ExpressionType); et {
	case "sql":
		return NewDefaultStorageQueryer(cfg, cloud), nil
	default:
		return nil, xerrors.New("Expression type you chose does not match. you must choose from 'SQL'")
	}
}

func NewDefaultStorageQueryer(cfg config.S3SelectConfig, cloud aws.Cloud) defaultStorageQueryer {
	csvInput := csvInputSerialization{
		fileHeaderInfo:   cfg.InputFileHeaderInfo,
		fieldDelimiter:   cfg.InputFieldDelimiter,
		compresssionType: cfg.InputCompressionType,
	}
	csvOutput := csvOutputSerialization{
		recordDelimiter: getOrDefault(cfg.OutputRecordDelimiter, cfg.InputFieldDelimiter),
	}
	return defaultStorageQueryer{
		cloud:     cloud,
		csvInput:  csvInput,
		csvOutput: csvOutput,
	}
}

var _ StorageQueryer = &defaultStorageQueryer{}

type defaultStorageQueryer struct {
	cloud      aws.Cloud
	logger     zap.Logger
	expression string
	csvInput   csvInputSerialization
	csvOutput  csvOutputSerialization
}

type csvInputSerialization struct {
	fileHeaderInfo   string
	compresssionType string
	fieldDelimiter   string
}

type csvOutputSerialization struct {
	recordDelimiter string
}

func getOrDefault(val, defaultVal string) string {
	if val != "" {
		return val
	}
	return defaultVal
}

func (s defaultStorageQueryer) Query(ctx context.Context, bucketName, objectKey string, writer io.Writer) {
	input := &s3sdk.SelectObjectContentInput{
		Bucket:          awssdk.String(bucketName),
		Key:             awssdk.String(objectKey),
		ExpressionType:  awssdk.String(s3.ExpressionTypeSql),
		Expression:      awssdk.String(s.expression),
		RequestProgress: &s3.RequestProgress{},
		InputSerialization: &s3.InputSerialization{
			CompressionType: awssdk.String(s.csvInput.compresssionType),
			CSV: &s3.CSVInput{
				FileHeaderInfo: awssdk.String(s.csvInput.fileHeaderInfo),
				FieldDelimiter: awssdk.String(s.csvInput.fieldDelimiter),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			CSV: &s3.CSVOutput{
				FieldDelimiter: awssdk.String(s.csvOutput.recordDelimiter),
			},
		},
	}
	resp, err := s.cloud.S3().SelectObjectContent(input)
	if err != nil {
		s.logger.Error("failed to execute s3select api", zap.String("error", err.Error()))
	}
	defer resp.EventStream.Close()

	for event := range resp.EventStream.Events() {
		// If the event type is `records`, it fetch the data from the message.
		v, ok := event.(*s3.RecordsEvent)
		if ok {
			writer.Write(v.Payload)
		}
	}
}
