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
				S3Config: config.S3Config{
					BucketName: "foo",
					KeyPrefix:  "bar",
					Format:     "CSV",
					SQL:        "select * from s3object limit 100",
				},
			},
			wantErr: nil,
		},
		{
			name: "bucket name is not specified",
			cfg: config.CliConfig{
				S3Config: config.S3Config{
					BucketName: "",
					KeyPrefix:  "hoge",
					Format:     "CSV",
					SQL:        "select * from s3object limit 100",
				},
			},
			wantErr: xerrors.New("bucket name must be specified"),
		},
		{
			name: "sql expression is not specified",
			cfg: config.CliConfig{
				S3Config: config.S3Config{
					BucketName: "foo",
					KeyPrefix:  "bar",
					Format:     "CSV",
					SQL:        "",
				},
			},
			wantErr: xerrors.New("sql must be specified"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cfg.S3Config.Validate()
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
