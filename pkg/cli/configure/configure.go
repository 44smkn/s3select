package configure

import (
	"fmt"
	"sort"

	"github.com/44smkn/s3select/pkg/cliutil"
	"github.com/44smkn/s3select/pkg/config"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

type ConfigureOptions struct {
	ProfileKey string
	ProfileVal *config.Profile

	Config func() (config.Config, error)
}

func NewCmdCongigure(f *cliutil.Factory) *cobra.Command {
	opts := ConfigureOptions{
		Config: f.Config,
	}
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Confugure s3select settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return configureRun(&opts)
		},
	}

	cmd.Flags().StringVarP(&opts.ProfileKey, "profile", "p", "default", "PROFILE")
	return cmd
}

func configureRun(opts *ConfigureOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}
	profiles, err := cfg.Profiles()
	if err != nil {
		return err
	}
	profile, ok := profiles[opts.ProfileKey]
	if !ok {
		xerrors.New("Not found your specified profile")
	}
	configPrompt(&profile)

	cfg.SetProfile(opts.ProfileKey, profile)
	return cfg.Write(config.ConfigFile())
}

func configPrompt(current *config.Profile) *config.Profile {
	current.InputSerialization.CompressionType =
		compressionTypePrompt(current.InputSerialization.CompressionType)

	current.InputSerialization.CSV.FieldDelimiter =
		inputCSVFieldDelimiterPrompt(current.InputSerialization.CSV.FieldDelimiter)

	return current
}

func selectPrompt(label string, items []string) string {
	p := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, _ := p.Run()
	return result
}

func compressionTypePrompt(currentVal string) string {
	label := fmt.Sprintf("Select CompressionType [%s]", currentVal)
	return selectPrompt(label, []string{
		s3.CompressionTypeNone,
		s3.CompressionTypeGzip,
		s3.CompressionTypeBzip2,
	})
}

var seprateCharacterMap = map[string]string{
	"SPACE": " ",
	"TAB": "	",
	"COMMMA":    ",",
	"SEMICOLON": ";",
}

func inputCSVFieldDelimiterPrompt(currentVal string) string {
	label := fmt.Sprintf("Select Input CSV Field Delimiter [%s]", getKeyByValue(seprateCharacterMap, currentVal))
	key := selectPrompt(label, KeySet(seprateCharacterMap))
	val, _ := seprateCharacterMap[key]
	return val
}

func KeySet(m map[string]string) []string {
	keyset := make([]string, 0, len(m))
	for k, _ := range m {
		keyset = append(keyset, k)
	}
	sort.Strings(keyset)
	return keyset
}

func getKeyByValue(m map[string]string, val string) string {
	for k, v := range m {
		if v == val {
			return k
		}
	}
	return ""
}
