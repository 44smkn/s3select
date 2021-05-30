package config

import (
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

const (
	flagBucketName             = "bucket-name"
	flagKeyPrefix              = "key-prefix"
	flagFormat                 = "format"
	flagExpressionType         = "expression-type"
	flagExpression             = "expression"
	flagInputCompressionType   = "input-compression-type"
	flagInputFileHeaderInfo    = "input-file-header-info"
	flagInputFieldDelimiter    = "input-field-delimiter"
	flagOutputFieldDelimiter   = "output-field-delimiter"
	defaultKeyPrefix           = "/"
	defaultFormat              = "CSV"
	defaultExpressionType      = "SQL"
	defaultInputFileHeaderInfo = "None"
	defaultInputFieldDelimiter = " "
)

type S3SelectConfig struct {
	BucketName            string
	KeyPrefix             string
	Format                string
	ExpressionType        string
	Expression            string
	InputCompressionType  string
	InputFileHeaderInfo   string
	InputFieldDelimiter   string
	OutputRecordDelimiter string
}

// BindFlags コマンドライン引数を構造体にバインドします
func (cfg *S3SelectConfig) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cfg.BucketName, flagBucketName, "", "The bucket name containing the object")
	fs.StringVar(&cfg.KeyPrefix, flagKeyPrefix, defaultKeyPrefix, "Key of the object to SELECT")
	fs.StringVar(&cfg.Format, flagFormat, defaultFormat, "Describes the format of the data in the object that is being queried")
	fs.StringVar(&cfg.Format, flagExpressionType, defaultExpressionType, "The type of the provided expression")
	fs.StringVar(&cfg.Expression, flagExpression, "", "The expression that is used to query the object")
	fs.StringVar(&cfg.InputCompressionType, flagInputCompressionType, "", "object's compression format. Valid values: NONE, GZIP, BZIP2")
	fs.StringVar(&cfg.InputFileHeaderInfo, flagInputFileHeaderInfo, defaultInputFileHeaderInfo, "Describes the first line of input. Valid values: None, IGNORE, USE")
	fs.StringVar(&cfg.InputFieldDelimiter, flagInputFieldDelimiter, defaultInputFieldDelimiter, "A single character used to separate individual fields in a record")
	fs.StringVar(&cfg.OutputRecordDelimiter, flagOutputFieldDelimiter, "", "A single character used to separate individual records in the output")
}

// Validate は設定内容に不備がないかを確認します
func (cfg *S3SelectConfig) Validate() error {
	if len(cfg.BucketName) == 0 {
		return xerrors.New("bucket name must be specified")
	}
	if len(cfg.Expression) == 0 {
		return xerrors.New("sql must be specified")
	}
	return nil
}
