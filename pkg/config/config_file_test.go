package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/44smkn/s3select/pkg/config"
)

func stubConfig(t *testing.T, configContents string) func() {
	t.Helper()
	original := config.ReadConfigFile
	config.ReadConfigFile = func(filename string) ([]byte, error) {
		switch filepath.Base(filename) {
		case "config.yaml":
			if configContents == "" {
				return []byte(nil), os.ErrNotExist
			} else {
				return []byte(configContents), nil
			}
		default:
			return []byte(nil), fmt.Errorf("read from unstubbed file: %q", filename)
		}
	}
	return func() {
		config.ReadConfigFile = original
	}
}
