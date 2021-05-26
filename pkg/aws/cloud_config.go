package aws

import "github.com/spf13/pflag"

const (
	flagAWSRegion = "region"
	defaultRegion = "ap-northeast-1"
)

type CloudConfig struct {
	// AWS Region for the S3 Bucket
	Region string
}

func (cfg *CloudConfig) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&cfg.Region, flagAWSRegion, defaultRegion, "AWS Region for the S3 Bucket")
}
