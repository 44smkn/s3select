package aws

import (
	"github.com/44smkn/s3select/pkg/aws/services"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Cloud interface {
	// S3 provides API to AWS S3
	S3() services.S3
}

func NewCloud(cfg CloudConfig) (Cloud, error) {
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(cfg.Region)))
	injectUserAgent(&sess.Handlers)
	return &defaultCloud{
		cfg: cfg,
		s3:  services.NewS3(sess),
	}, nil
}

var _ Cloud = &defaultCloud{}

type defaultCloud struct {
	cfg CloudConfig
	s3  services.S3
}

func (c *defaultCloud) S3() services.S3 {
	return c.s3
}
