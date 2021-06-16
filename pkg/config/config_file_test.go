package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/44smkn/s3select/pkg/config"
	"github.com/aws/aws-sdk-go/service/s3"
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
      csvInput:
        fieldDelimiter: ';'
        quoteCharacter: '"'
      compressionType: BZIP2
      jsonInput:
        type: DOCUMENT
    outputSerialization:
      csvInput:
        fieldDelimiter: ' '
        quoteCharacter: '"'
      jsonInput:
        recordDelimiter: \n
  albAcesssLog:
    expressionType: SQL
    inputSerialization:
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
				ExpressionType: s3.ExpressionTypeSql,
				InputSerialization: config.InputSerialization{
					CompressionType: s3.CompressionTypeBzip2,
					CSV: config.CSVInput{
						FieldDelimiter: ";",
						QuoteCharacter: `"`,
					},
					JSON: config.JSONInput{
						Type: s3.JSONTypeDocument,
					},
				},
				OutputSerialization: config.OutputSerialization{
					CSV: config.CSVOutput{
						FieldDelimiter: " ",
						QuoteCharacter: `"`,
					},
					JSON: config.JSONOutput{
						RecordDelimiter: `\n`,
					},
				},
			},
			"albAcesssLog": {
				ExpressionType: s3.ExpressionTypeSql,
				InputSerialization: config.InputSerialization{
					CompressionType: s3.CompressionTypeGzip,
					CSV: config.CSVInput{
						AllowQuotedRecordDelimiter: true,
						Comments:                   "#",
						FieldDelimiter:             " ",
						FileHeaderInfo:             s3.FileHeaderInfoIgnore,
						QuoteCharacter:             `"`,
						QuoteEscapeCharacter:       `"`,
						RecordDelimiter:            `\n`,
					},
					JSON: config.JSONInput{
						Type: s3.JSONTypeDocument,
					},
				},
				OutputSerialization: config.OutputSerialization{
					CSV: config.CSVOutput{
						FieldDelimiter:       ",",
						QuoteCharacter:       `"`,
						QuoteEscapeCharacter: `"`,
						QuoteFields:          s3.QuoteFieldsAsneeded,
						RecordDelimiter:      `\n`,
					},
					JSON: config.JSONOutput{
						RecordDelimiter: `\n`,
					},
				},
			},
		},
	}
	if diff := cmp.Diff(cfg, want); diff != "" {
		t.Errorf("ParseDefaultConfig() mismatch (-want +got):\n%s", diff)
	}

}
