package config

import (
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

const (
	flagBucketName           = "bucket-name"
	flagKeyPrefix            = "key-prefix"
	flagFormat               = "format"
	flagSQLExpression        = "sql-expression"
	flagInputCompressionType = "input-compression-type"
	defaultKeyPrefix         = "/"
	defaultFormat            = "CSV"
)

type S3SelectConfig struct {
	BucketName           string
	KeyPrefix            string
	Format               string
	SQLExpression        string
	InputCompressionType string
}

// BindFlags コマンドライン引数を構造体にバインドします
func (cfg *S3SelectConfig) BindFlags(fs *pflag.FlagSet) {
	// required
	fs.StringVar(&cfg.BucketName, flagBucketName, "", "The bucket name containing the object")
	fs.StringVar(&cfg.KeyPrefix, flagKeyPrefix, defaultKeyPrefix, "Key of the object to SELECT")
	fs.StringVar(&cfg.Format, flagFormat, defaultFormat, "Describes the format of the data in the object that is being queried")
	fs.StringVar(&cfg.SQLExpression, flagSQLExpression, "", "The expression that is used to query the object")
	// option
	fs.StringVar(&cfg.InputCompressionType, flagInputCompressionType, "", "object's compression format. Valid values: NONE, GZIP, BZIP2")
}

// Validate は設定内容に不備がないかを確認します
func (cfg *S3SelectConfig) Validate() error {
	if len(cfg.BucketName) == 0 {
		return xerrors.New("bucket name must be specified")
	}
	if len(cfg.SQLExpression) == 0 {
		return xerrors.New("sql must be specified")
	}
	return nil
}
