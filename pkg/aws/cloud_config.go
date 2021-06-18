package aws

const (
	flagAWSRegion = "region"
	defaultRegion = "ap-northeast-1"
)

type CloudConfig struct {
	// AWS Region for the S3 Bucket
	Region string
}
