package config

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	S3SELECT_CONFIG_DIR = "S3SELECT_CONFIG_DIR"
)

type FileConfig struct {
	AWSRegion string   `yaml:"awsRegion"`
	Proflies  Profiles `yaml:"profiles"`
}

type Profiles = map[string]Profile

type Profile struct {
	ExpressionType      string              `yaml:"expressionType"`
	InputSerialization  InputSerialization  `yaml:"inputSerialization"`
	OutputSerialization OutputSerialization `yaml:"outputSerialization"`
}

type InputSerialization struct {
	CSV             CSVInput     `yaml:"csvInput"`
	CompressionType string       `yaml:"compressionType"`
	JSON            JSONInput    `yaml:"jsonInput"`
	Parquet         ParquetInput `yaml:"parquetInput"`
}

type CSVInput struct {
	AllowQuotedRecordDelimiter bool   `yaml:"allowQuotedRecordDelimiter"`
	Comments                   string `yaml:"comments"`
	FieldDelimiter             string `yaml:"fieldDelimiter"`
	FileHeaderInfo             string `yaml:"fileHeaderInfo"`
	QuoteCharacter             string `yaml:"quoteCharacter"`
	QuoteEscapeCharacter       string `yaml:"quoteEscapeCharacter"`
	RecordDelimiter            string `yaml:"recordDelimiter"`
}

type JSONInput struct {
	Type string `yaml:"type"`
}

type ParquetInput struct{}

type OutputSerialization struct {
	CSV  CSVOutput  `yaml:"csvInput"`
	JSON JSONOutput `yaml:"jsonInput"`
}

type CSVOutput struct {
	FieldDelimiter       string `yaml:"fieldDelimiter"`
	QuoteCharacter       string `yaml:"quoteCharacter"`
	QuoteEscapeCharacter string `yaml:"quoteEscapeCharacter"`
	QuoteFields          string `yaml:"quoteFields"`
	RecordDelimiter      string `yaml:"recordDelimiter"`
}

type JSONOutput struct {
	RecordDelimiter string `yaml:"recordDelimiter"`
}

func (c *FileConfig) GetAWSRegion() string {
	return c.AWSRegion
}

func (c *FileConfig) Profiles() (Profiles, error) {
	return c.Proflies, nil
}

func ConfigDir() string {
	if path := os.Getenv(S3SELECT_CONFIG_DIR); path != "" {
		return path
	}
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "s3select")
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
	cfg := &FileConfig{} // TODO: Set default Value
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
