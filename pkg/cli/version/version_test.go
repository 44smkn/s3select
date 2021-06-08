package version_test

import (
	"testing"

	"github.com/44smkn/s3select/pkg/cli/version"
)

func TestFormat(t *testing.T) {
	expects := "s3select version 0.1.0 (2021-06-02)\n"
	if got := version.Format("0.1.0", "2021-06-02"); got != expects {
		t.Errorf("Format() = %q, wants %q", got, expects)
	}
}
