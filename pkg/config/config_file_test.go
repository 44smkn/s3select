package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/44smkn/s3select/pkg/config"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/go-cmp/cmp"
)

func stubConfig(t *testing.T, configContents string) func() {
	t.Helper()
	original := config.ReadConfigFile
	config.ReadConfigFile = func(filename string) ([]byte, error) {
		switch filepath.Base(filename) {
		case "config.yaml":
			if configContents == "" {
				return []byte(nil), os.ErrNotExist
			} else {
				return []byte(configContents), nil
			}
		default:
			return []byte(nil), fmt.Errorf("read from unstubbed file: %q", filename)
		}
	}
	return func() {
		config.ReadConfigFile = original
	}
}

func TestParseDefaultConfig(t *testing.T) {
	defer stubConfig(t, `
awsRegion: ap-northeast-1
profiles:
  default:
    expressionType: SQL
    inputSerialization:
      formatType: csv
      csvInput:
        fieldDelimiter: ';'
        quoteCharacter: '"'
      compressionType: BZIP2
      jsonInput:
        type: DOCUMENT
    outputSerialization:
      formatType: csv
      csvInput:
        fieldDelimiter: ' '
        quoteCharacter: '"'
      jsonInput:
        recordDelimiter: \n
  albAcesssLog:
    expressionType: SQL
    inputSerialization:
      formatType: csv
      csvInput:
        allowQuotedRecordDelimiter: true
        comments: '#'
        fieldDelimiter: ' '
        fileHeaderInfo: IGNORE
        quoteCharacter: '"'
        quoteEscapeCharacter: '"'
        recordDelimiter: \n
      compressionType: GZIP
      jsonInput:
        type: DOCUMENT
    outputSerialization:
      formatType: csv
      csvInput:
        fieldDelimiter: ','
        quoteCharacter: '"'
        quoteEscapeCharacter: '"'
        quoteFields: ASNEEDED
        recordDelimiter: \n
      jsonInput:
        recordDelimiter: \n
`)()
	cfg, err := config.ParseDefaultConfig()
	if err != nil {
		t.Errorf("ParseDefaultConfig() error: %v", err)
	}
	want := &config.FileConfig{
		AWSRegion: "ap-northeast-1",
		Proflies: map[string]config.Profile{
			"default": {
				ExpressionType: s3sdk.ExpressionTypeSql,
				InputSerialization: &config.InputSerialization{
					FormatType:      "csv",
					CompressionType: awssdk.String(s3sdk.CompressionTypeBzip2),
					CSV: &config.CSVInput{
						FieldDelimiter: awssdk.String(";"),
						QuoteCharacter: awssdk.String(`"`),
					},
					JSON: &config.JSONInput{
						Type: awssdk.String(s3sdk.JSONTypeDocument),
					},
				},
				OutputSerialization: &config.OutputSerialization{
					FormatType: "csv",
					CSV: &config.CSVOutput{
						FieldDelimiter: awssdk.String(" "),
						QuoteCharacter: awssdk.String(`"`),
					},
					JSON: &config.JSONOutput{
						RecordDelimiter: awssdk.String(`\n`),
					},
				},
			},
			"albAcesssLog": {
				ExpressionType: s3.ExpressionTypeSql,
				InputSerialization: &config.InputSerialization{
					FormatType:      "csv",
					CompressionType: awssdk.String(s3.CompressionTypeGzip),
					CSV: &config.CSVInput{
						AllowQuotedRecordDelimiter: awssdk.Bool(true),
						Comments:                   awssdk.String("#"),
						FieldDelimiter:             awssdk.String(" "),
						FileHeaderInfo:             awssdk.String(s3.FileHeaderInfoIgnore),
						QuoteCharacter:             awssdk.String(`"`),
						QuoteEscapeCharacter:       awssdk.String(`"`),
						RecordDelimiter:            awssdk.String(`\n`),
					},
					JSON: &config.JSONInput{
						Type: awssdk.String(s3.JSONTypeDocument),
					},
				},
				OutputSerialization: &config.OutputSerialization{
					FormatType: "csv",
					CSV: &config.CSVOutput{
						FieldDelimiter:       awssdk.String(","),
						QuoteCharacter:       awssdk.String(`"`),
						QuoteEscapeCharacter: awssdk.String(`"`),
						QuoteFields:          awssdk.String(s3.QuoteFieldsAsneeded),
						RecordDelimiter:      awssdk.String(`\n`),
					},
					JSON: &config.JSONOutput{
						RecordDelimiter: awssdk.String(`\n`),
					},
				},
			},
		},
	}
	if diff := cmp.Diff(cfg, want); diff != "" {
		t.Errorf("ParseDefaultConfig() mismatch (-want +got):\n%s", diff)
	}

}
