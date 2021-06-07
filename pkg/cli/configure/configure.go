package configure

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

type ConfigureOptions struct {
	Profile string
}

func NewCmdCongigure() *cobra.Command {
	opts := ConfigureOptions{}
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Confugure s3select settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return configureRun(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Profile, "profile", "p", "default", "PROFILE")
	return cmd
}

func configureRun(opts *ConfigureOptions) error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		xerrors.Errorf("failed to obtain user's home directory: %w", err)
	}

	s3selectDir := filepath.Join(dirname, ".s3select")
	if _, err := os.Stat(s3selectDir); os.IsNotExist(err) {
		if err := os.Mkdir(s3selectDir, 0700); err != nil {
			return xerrors.Errorf("failed to make directory '%s': %w", s3selectDir, err)
		}
	}

	s3selectCfgFile := filepath.Join(s3selectDir, "config.yaml")
	if _, err := os.Stat(s3selectCfgFile); os.IsNotExist(err) {
		if err := os.WriteFile(s3selectCfgFile, []byte(""), 0600); err != nil {
			return xerrors.Errorf("failed to make config file '%s': %w", s3selectDir, err)
		}
	}

	return nil
}
