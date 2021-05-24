package config

import (
	"errors"

	"github.com/spf13/pflag"
)

const (
	flagBucketName   = "bucket-name"
	flagKeyPrefix    = "key-prefix"
	flagFormat       = "format"
	flagSQL          = "sql"
	defaultKeyPrefix = "/"
	defaultFormat    = "CSV"
)

type FlagConfig struct {
	BucketName string
	KeyPrefix  string
	Format     string
	SQL        string
}

// BindFlags コマンドライン引数を構造体にバインドします
func (cfg *FlagConfig) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cfg.BucketName, flagBucketName, "", "The bucket name containing the object")
	fs.StringVar(&cfg.KeyPrefix, flagKeyPrefix, defaultKeyPrefix, "Key of the object to SELECT")
	fs.StringVar(&cfg.Format, flagFormat, defaultFormat, "Describes the format of the data in the object that is being queried")
	fs.StringVar(&cfg.SQL, flagSQL, "", "The expression that is used to query the object")
}

// Validate は設定内容に不備がないかを確認します
func (cfg *FlagConfig) Validate() error {
	if len(cfg.BucketName) == 0 {
		return errors.New("bucket name must be specified")
	}
	if len(cfg.SQL) == 0 {
		return errors.New("sql must be specified")
	}
	return nil
}
