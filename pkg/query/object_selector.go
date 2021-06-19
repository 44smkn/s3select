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
	"golang.org/x/xerrors"
)

type ObjectSelector interface {
	Select(context.Context, *ObjectMetadata, string, io.Writer) error
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

func (s defaultObjectSelector) Select(ctx context.Context, meta *ObjectMetadata, expression string, writer io.Writer) error {
	is, err := newInputSerialization(s.cfg.InputSerialization)
	if err != nil {
		return err
	}
	os, err := newOutputSerialization(s.cfg.OutputSerialization)
	if err != nil {
		return err
	}
	input := &s3sdk.SelectObjectContentInput{
		Bucket:              awssdk.String(meta.BucketName),
		Key:                 awssdk.String(meta.ObjectKey),
		ExpressionType:      awssdk.String(s.cfg.ExpressionType),
		Expression:          awssdk.String(expression),
		RequestProgress:     &s3.RequestProgress{},
		InputSerialization:  is,
		OutputSerialization: os,
	}
	resp, err := s.cloud.S3().SelectObjectContentWithContext(ctx, input)
	if err != nil {
		return xerrors.Errorf("failed to execute s3api: %w", err)
	}
	defer resp.EventStream.Close()

	for event := range resp.EventStream.Events() {
		// If the event type is `records`, it fetch the data from the message.
		v, ok := event.(*s3.RecordsEvent)
		if ok {
			writer.Write(v.Payload)
		}
	}
	return nil
}

func newInputSerialization(cfg config.InputSerialization) (*s3sdk.InputSerialization, error) {
	switch cfg.FormatType {
	case config.S3SELECT_INPUT_FORMAT_CSV:
		return &s3.InputSerialization{
			CompressionType: cfg.CompressionType,
			CSV: &s3.CSVInput{
				AllowQuotedRecordDelimiter: cfg.CSV.AllowQuotedRecordDelimiter,
				Comments:                   cfg.CSV.Comments,
				FieldDelimiter:             cfg.CSV.FieldDelimiter,
				FileHeaderInfo:             cfg.CSV.FileHeaderInfo,
				QuoteCharacter:             cfg.CSV.QuoteCharacter,
				QuoteEscapeCharacter:       cfg.CSV.QuoteEscapeCharacter,
				RecordDelimiter:            cfg.CSV.RecordDelimiter,
			},
		}, nil
	case config.S3SELECT_INPUT_FORMAT_JSON:
		return &s3sdk.InputSerialization{
			JSON: &s3.JSONInput{
				Type: cfg.JSON.Type,
			},
		}, nil
	default:
		return nil, xerrors.New("choose a input format type from: [json, csv]")
	}
}

func newOutputSerialization(cfg config.OutputSerialization) (*s3sdk.OutputSerialization, error) {
	switch cfg.FormatType {
	case config.S3SELECT_INPUT_FORMAT_CSV:
		return &s3.OutputSerialization{
			CSV: &s3.CSVOutput{
				FieldDelimiter:       cfg.CSV.FieldDelimiter,
				QuoteCharacter:       cfg.CSV.QuoteCharacter,
				QuoteEscapeCharacter: cfg.CSV.QuoteEscapeCharacter,
				QuoteFields:          cfg.CSV.QuoteFields,
				RecordDelimiter:      cfg.CSV.RecordDelimiter,
			},
		}, nil
	case config.S3SELECT_INPUT_FORMAT_JSON:
		return &s3.OutputSerialization{
			JSON: &s3.JSONOutput{
				RecordDelimiter: cfg.JSON.RecordDelimiter,
			},
		}, nil
	default:
		return nil, xerrors.New("choose a output format type from: [json, csv]")
	}
}
