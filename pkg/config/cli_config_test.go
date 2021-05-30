package config_test

import (
	"testing"

	"github.com/44smkn/s3selecgo/pkg/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		cfg     config.CliConfig
		wantErr error
	}{
		{
			name: "normal case",
			cfg: config.CliConfig{
				S3SelectConfig: config.S3SelectConfig{
					BucketName: "foo",
					KeyPrefix:  "bar",
					Format:     "CSV",
					Expression: "select * from s3object limit 100",
				},
			},
			wantErr: nil,
		},
		{
			name: "bucket name is not specified",
			cfg: config.CliConfig{
				S3SelectConfig: config.S3SelectConfig{
					BucketName: "",
					KeyPrefix:  "hoge",
					Format:     "CSV",
					Expression: "select * from s3object limit 100",
				},
			},
			wantErr: xerrors.New("bucket name must be specified"),
		},
		{
			name: "sql expression is not specified",
			cfg: config.CliConfig{
				S3SelectConfig: config.S3SelectConfig{
					BucketName: "foo",
					KeyPrefix:  "bar",
					Format:     "CSV",
					Expression: "",
				},
			},
			wantErr: xerrors.New("sql must be specified"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cfg.S3SelectConfig.Validate()
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
