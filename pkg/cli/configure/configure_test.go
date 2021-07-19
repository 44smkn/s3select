package configure_test

import (
	"testing"

	"github.com/44smkn/s3select/pkg/cli"
	"github.com/44smkn/s3select/pkg/cli/configure"
)

func TestRegionPrompt(t *testing.T) {
	tests := []struct {
		name    string
		current string
		input   string
		want    string
	}{
		{
			name:    "test1",
			current: "ap-northeast-1",
			input:   "us-east-1",
			want:    "us-east-1",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer cli.StubSurveyAskOne(tt.input)()
			if got := configure.RegionPrompt(tt.current); got != tt.want {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}
