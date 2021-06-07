package configure

import (
	"context"
	"fmt"
	"os"

	"github.com/44smkn/s3selecgo/pkg/aws"
	"github.com/44smkn/s3selecgo/pkg/config"
	"github.com/44smkn/s3selecgo/pkg/query"
	awssdk "github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

type ConfigureOptions struct {
	Profile    string
	BucketName string
	KeyPrefix  string
}

func NewCmdQuery() *cobra.Command {
	opts := ConfigureOptions{}
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Execute S3 SELECT",
		RunE: func(cmd *cobra.Command, args []string) error {
			return queryRun(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Profile, "profile", "p", "default", "PROFILE")
	return cmd
}

func queryRun(opts *ConfigureOptions) error {
	awsCfg := aws.CloudConfig{
		// TODO: read profile
		Region: "ap-northeast-1",
	}
	cloud, err := aws.NewCloud(awsCfg)
	if err != nil {
		return err
	}

	ctx := context.Background()

	req := s3sdk.ListObjectsV2Input{
		Bucket: awssdk.String(opts.BucketName),
		Prefix: awssdk.String(opts.KeyPrefix),
	}
	objects, err := cloud.S3().ListObjectsV2AsList(ctx, &req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute s3listobeject api: %s", err.Error())
		return err
	}

	s3selectCfg := config.S3SelectConfig{}
	sq := query.NewDefaultStorageQueryer(s3selectCfg, cloud)
	for _, o := range objects {
		sq.Query(ctx, opts.BucketName, *o.Key, os.Stdout)
	}

	return nil
}
