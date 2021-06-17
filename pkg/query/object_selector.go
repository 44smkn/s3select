package query

import (
	"context"
	"io"

	"github.com/44smkn/s3select/pkg/aws"
	"github.com/44smkn/s3select/pkg/config"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

type ObjectSelector interface {
	Select(context.Context, *ObjectMetadata, string, io.Writer)
}

func NewDefaultObjectSelector(profile *config.Profile, cloud aws.Cloud, logger *zap.Logger) defaultObjectSelector {
	return defaultObjectSelector{
		cfg:    profile,
		cloud:  cloud,
		logger: logger,
	}
}

var _ ObjectSelector = &defaultObjectSelector{}

type defaultObjectSelector struct {
	cfg    *config.Profile
	cloud  aws.Cloud
	logger *zap.Logger
}

type ObjectMetadata struct {
	BucketName string
	ObjectKey  string
}

func (s defaultObjectSelector) Select(ctx context.Context, meta *ObjectMetadata, expression string, writer io.Writer) {
	input := &s3sdk.SelectObjectContentInput{
		Bucket:          awssdk.String(meta.BucketName),
		Key:             awssdk.String(meta.ObjectKey),
		ExpressionType:  awssdk.String(s.cfg.ExpressionType),
		Expression:      awssdk.String(expression),
		RequestProgress: &s3.RequestProgress{},
		InputSerialization: &s3.InputSerialization{
			CompressionType: awssdk.String(s.cfg.ExpressionType),
			CSV: &s3.CSVInput{
				AllowQuotedRecordDelimiter: awssdk.Bool(s.cfg.InputSerialization.CSV.AllowQuotedRecordDelimiter),
				Comments:                   awssdk.String(s.cfg.InputSerialization.CSV.Comments),
				FieldDelimiter:             awssdk.String(s.cfg.InputSerialization.CSV.FieldDelimiter),
				FileHeaderInfo:             awssdk.String(s.cfg.InputSerialization.CSV.FileHeaderInfo),
				QuoteCharacter:             awssdk.String(s.cfg.InputSerialization.CSV.QuoteCharacter),
				QuoteEscapeCharacter:       awssdk.String(s.cfg.InputSerialization.CSV.QuoteEscapeCharacter),
				RecordDelimiter:            awssdk.String(s.cfg.InputSerialization.CSV.RecordDelimiter),
			},
			JSON: &s3.JSONInput{
				Type: awssdk.String(s.cfg.InputSerialization.JSON.Type),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			CSV: &s3.CSVOutput{
				FieldDelimiter:       awssdk.String(s.cfg.OutputSerialization.CSV.FieldDelimiter),
				QuoteCharacter:       awssdk.String(s.cfg.OutputSerialization.CSV.QuoteCharacter),
				QuoteEscapeCharacter: awssdk.String(s.cfg.OutputSerialization.CSV.QuoteEscapeCharacter),
				QuoteFields:          awssdk.String(s.cfg.OutputSerialization.CSV.QuoteFields),
				RecordDelimiter:      awssdk.String(s.cfg.OutputSerialization.CSV.RecordDelimiter),
			},
			JSON: &s3sdk.JSONOutput{
				RecordDelimiter: awssdk.String(s.cfg.OutputSerialization.JSON.RecordDelimiter),
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
