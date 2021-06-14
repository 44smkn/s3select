package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
)

const (
	S3SELECT_CONFIG_DIR = "S3SELECT_CONFIG_DIR"
)

type FileConfig struct {
	AWSRegion string   `yaml:"awsRegion,omitempty"`
	Proflies  Profiles `yaml:"profiles,omitempty"`
}

type Profiles = map[string]Profile

type Profile struct {
	ExpressionType      string              `yaml:"expressionType,omitempty"`
	InputSerialization  InputSerialization  `yaml:"inputSerialization,omitempty"`
	OutputSerialization OutputSerialization `yaml:"outputSerialization,omitempty"`
}

type InputSerialization struct {
	CSV             CSVInput     `yaml:"csvInput,omitempty"`
	CompressionType string       `yaml:"compressionType,omitempty"`
	JSON            JSONInput    `yaml:"jsonInput,omitempty"`
	Parquet         ParquetInput `yaml:"parquetInput,omitempty"`
}

type CSVInput struct {
	AllowQuotedRecordDelimiter bool   `yaml:"allowQuotedRecordDelimiter,omitempty"`
	Comments                   string `yaml:"comments,omitempty"`
	FieldDelimiter             string `yaml:"fieldDelimiter,omitempty"`
	FileHeaderInfo             string `yaml:"fileHeaderInfo,omitempty"`
	QuoteCharacter             string `yaml:"quoteCharacter,omitempty"`
	QuoteEscapeCharacter       string `yaml:"quoteEscapeCharacter,omitempty"`
	RecordDelimiter            string `yaml:"recordDelimiter,omitempty"`
}

type JSONInput struct {
	Type string `yaml:"type,omitempty"`
}

type ParquetInput struct{}

type OutputSerialization struct {
	CSV  CSVOutput  `yaml:"csvInput,omitempty"`
	JSON JSONOutput `yaml:"jsonInput,omitempty"`
}

type CSVOutput struct {
	FieldDelimiter       string `yaml:"fieldDelimiter,omitempty"`
	QuoteCharacter       string `yaml:"quoteCharacter,omitempty"`
	QuoteEscapeCharacter string `yaml:"quoteEscapeCharacter,omitempty"`
	QuoteFields          string `yaml:"quoteFields,omitempty"`
	RecordDelimiter      string `yaml:"recordDelimiter,omitempty"`
}

type JSONOutput struct {
	RecordDelimiter string `yaml:"recordDelimiter,omitempty"`
}

func (c *FileConfig) GetAWSRegion() string {
	return c.AWSRegion
}

func (c *FileConfig) Profiles() (Profiles, error) {
	return c.Proflies, nil
}

func (c *FileConfig) SetProfile(key string, newProfile Profile) {
	c.Proflies[key] = newProfile
}

func ConfigDir() string {
	if path := os.Getenv(S3SELECT_CONFIG_DIR); path != "" {
		return path
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "s3select")
}

func ConfigFile() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

func ParseDefaultConfig() (Config, error) {
	return parseConfig(ConfigFile())
}

func parseConfig(filename string) (Config, error) {
	if !fileExists(filename) {
		initConfigFile(filename)
	}
	return parseConfigFile(filename)
}

func initConfigFile(filename string) {
	cfg := &FileConfig{
		AWSRegion: "us-west-2",
		Proflies: map[string]Profile{
			"default": {
				ExpressionType: s3.ExpressionTypeSql,
				InputSerialization: InputSerialization{
					CompressionType: s3.CompressionTypeNone,
					CSV: CSVInput{
						AllowQuotedRecordDelimiter: false,
						FieldDelimiter:             ",",
						QuoteCharacter:             `"`,
					},
					JSON: JSONInput{
						Type: s3.JSONTypeDocument,
					},
				},
				OutputSerialization: OutputSerialization{
					CSV: CSVOutput{
						FieldDelimiter: ",",
						QuoteCharacter: `"`,
					},
					JSON: JSONOutput{
						RecordDelimiter: "\n",
					},
				},
			},
		},
	}
	cfg.Write(filename)
}

func (c *FileConfig) Write(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return writeConfigFile(filename, data)
}

func writeConfigFile(filename string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(filename), 0771)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func parseConfigFile(filename string) (Config, error) {
	data, err := readConfigFile(filename)
	if err != nil {
		return nil, err
	}

	var config FileConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func readConfigFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func fileExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && !f.IsDir()
}