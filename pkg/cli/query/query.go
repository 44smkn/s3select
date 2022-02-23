package query

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/44smkn/s3select/pkg/aws"
	"github.com/44smkn/s3select/pkg/cli"
	"github.com/44smkn/s3select/pkg/config"
	"github.com/44smkn/s3select/pkg/query"
	awssdk "github.com/aws/aws-sdk-go/aws"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type QueryOptions struct {
	Out    io.Writer
	ErrOut io.Writer
	Logger *zap.Logger

	Profile    string
	BucketName string
	KeyPrefix  string
	Expression string
	AWSRegion  string

	Config func() (config.Config, error)
}

func NewCmdQuery(f *cli.Factory) *cobra.Command {
	opts := &QueryOptions{
		Out:    f.Out,
		ErrOut: f.ErrOut,
		Logger: f.Logger,
		Config: f.Config,
	}
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Execute S3 SELECT",
		RunE: func(cmd *cobra.Command, args []string) error {
			return queryRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Profile, "profile", "p", "default", "PROFILE")
	cmd.Flags().StringVarP(&opts.BucketName, "bucket", "b", "", "bucket name")
	cmd.Flags().StringVarP(&opts.KeyPrefix, "key-prefix", "k", "/", "key prefix")
	cmd.Flags().StringVarP(&opts.Expression, "expression", "e", "", "expression")
	cmd.Flags().StringVar(&opts.Expression, "aws-region", "", "expression")
	return cmd
}

func queryRun(opts *QueryOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	awsRegion := cfg.GetAWSRegion()
	if opts.AWSRegion != "" {
		awsRegion = opts.AWSRegion
	}
	awsCfg := aws.CloudConfig{
		Region: awsRegion,
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

	profiles, err := cfg.Profiles()
	if err != nil {
		return err
	}
	profile, ok := profiles[opts.Profile]
	if !ok {
		return xerrors.New("your specified profile is not found")
	}
	sq := query.NewDefaultObjectSelector(&profile, cloud, opts.Logger)
	for _, o := range objects {
		meta := &query.ObjectMetadata{
			BucketName: opts.BucketName,
			ObjectKey:  *o.Key,
		}
		err = sq.Select(ctx, meta, opts.Expression, opts.Out)
		if err != nil {
			return err
		}
	}

	return nil
}
