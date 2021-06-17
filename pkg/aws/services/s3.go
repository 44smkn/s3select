package services

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

//go:generate mockgen -destination=../../../mocks/aws/services/mock_s3.go -package=mock_services github.com/44smkn/s3select/pkg/aws/services S3
type S3 interface {
	s3iface.S3API

	// wrapper to ListObjectsV2PagesWithContext, which aggregates paged results into list.
	ListObjectsV2AsList(ctx context.Context, input *s3.ListObjectsV2Input) ([]*s3.Object, error)
}

func NewS3(session *session.Session) S3 {
	return &defaultS3{
		S3API: s3.New(session),
	}
}

// default implementation for S3.
type defaultS3 struct {
	s3iface.S3API
}

func (c *defaultS3) ListObjectsV2AsList(ctx context.Context, input *s3.ListObjectsV2Input) ([]*s3.Object, error) {
	var result []*s3.Object
	if err := c.ListObjectsV2PagesWithContext(ctx, input, func(output *s3.ListObjectsV2Output, _ bool) bool {
		result = append(result, output.Contents...)
		return true
	}); err != nil {
		return nil, err
	}
	return result, nil
}
