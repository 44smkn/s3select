package cli

import (
	"io"
	"os"

	"github.com/44smkn/s3select/pkg/config"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
)

type Factory struct {
	In     io.ReadCloser
	Out    io.Writer
	ErrOut io.Writer

	Config     func() (config.Config, error)
	Logger     *zap.Logger
	Executable string
}

func NewFactory(appVersion string, logger *zap.Logger) (*Factory, error) {
	var cachedCfg config.Config
	var cfgErr error
	cfgFunc := func() (config.Config, error) {
		if cachedCfg != nil || cfgErr != nil {
			return cachedCfg, cfgErr
		}
		cachedCfg, cfgErr = config.ParseDefaultConfig()
		return cachedCfg, cfgErr
	}

	executable := "s3select"
	if exe, err := os.Executable(); err == nil {
		executable = exe
	}

	return &Factory{
		In:     os.Stdin,
		Out:    colorable.NewColorable(os.Stdout),
		ErrOut: colorable.NewColorable(os.Stderr),

		Config:     cfgFunc,
		Logger:     logger,
		Executable: executable,
	}, nil
}
