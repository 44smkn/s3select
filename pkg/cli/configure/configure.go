package configure

import (
	"fmt"
	"sort"

	"github.com/44smkn/s3select/pkg/cli"
	"github.com/44smkn/s3select/pkg/config"
	"github.com/AlecAivazis/survey/v2"
	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

var (
	awsRegions = []string{
		"us-east-2",
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"af-south-1",
		"ap-east-1",
		"ap-south-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"ca-central-1",
		"cn-north-1",
		"cn-northwest-1",
		"eu-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-south-1",
		"eu-west-3",
		"eu-north-1",
		"me-south-1",
		"sa-east-1",
	}
)

type ConfigureOptions struct {
	ProfileKey      string
	ProfileVal      *config.Profile
	CongigureDetail bool

	Config func() (config.Config, error)
}

func NewCmdCongigure(f *cli.Factory) *cobra.Command {
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
	cmd.Flags().BoolVar(&opts.CongigureDetail, "detail", false, "configure s3select in detail")
	return cmd
}

func configureRun(opts *ConfigureOptions) error {
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	region, err := regionPrompt(cfg.GetAWSRegion())
	if err != nil {
		return err
	}
	cfg.SetAWSRegion(region)

	profiles, err := cfg.Profiles()
	if err != nil {
		return err
	}

	profile, ok := profiles[opts.ProfileKey]
	if !ok {
		profile = config.NewDefaultProfile()
	}

	is, err := inputSerializationPrompt(profile.InputSerialization, opts.CongigureDetail)
	if err != nil {
		return err
	}
	os, err := outputSerializationPrompt(profile.OutputSerialization, opts.CongigureDetail)
	if err != nil {
		return err
	}
	profile.SetSerializations(is, os)
	cfg.SetProfile(opts.ProfileKey, profile)
	return cfg.Write(config.ConfigFile())
}

func regionPrompt(defaultRegion string) (string, error) {
	region := defaultRegion
	prompt := &survey.Input{
		Message: "Region",
	}
	validator := func(input interface{}) error {
		for _, r := range awsRegions {
			if input == r {
				return nil
			}
		}
		return xerrors.New("your specified region does not exists")
	}
	err := cli.SurvayAskOne(prompt, &region, survey.WithValidator(validator))
	return region, err
}

func inputSerializationPrompt(cfg *config.InputSerialization, detail bool) (*config.InputSerialization, error) {
	inputFormat := inputFormatTypePrompt(cfg.FormatType)
	switch inputFormat {
	case config.S3SELECT_INPUT_FORMAT_CSV:
		if !detail {
			return &config.InputSerialization{
				FormatType:      config.S3SELECT_INPUT_FORMAT_CSV,
				CompressionType: awssdk.String(compressionTypePrompt(*cfg.CompressionType)),
				CSV: &config.CSVInput{
					FieldDelimiter: awssdk.String(csvFieldDelimiterPrompt(*cfg.CSV.FieldDelimiter)),
				},
			}, nil
		}
	}
	return nil, xerrors.New("choose a input format type from: [json, csv]")
}

func inputFormatTypePrompt(current string) string {
	label := fmt.Sprintf("Select Input Format Type [%s]", current)
	return selectPrompt(label, []string{
		config.S3SELECT_INPUT_FORMAT_CSV,
		config.S3SELECT_INPUT_FORMAT_JSON,
	})
}

func compressionTypePrompt(current string) string {
	label := fmt.Sprintf("Select CompressionType [%s]", current)
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

func csvFieldDelimiterPrompt(current string) string {
	label := fmt.Sprintf("Select Input CSV Field Delimiter [%s]", getKeyByValue(seprateCharacterMap, current))
	key := selectPrompt(label, KeySet(seprateCharacterMap))
	val, _ := seprateCharacterMap[key]
	return val
}

func outputSerializationPrompt(cfg *config.OutputSerialization, detail bool) (*config.OutputSerialization, error) {
	inputFormat := outputFormatTypePrompt(cfg.FormatType)
	switch inputFormat {
	case config.S3SELECT_INPUT_FORMAT_CSV:
		if !detail {
			return &config.OutputSerialization{
				FormatType: config.S3SELECT_INPUT_FORMAT_CSV,
				CSV: &config.CSVOutput{
					FieldDelimiter: awssdk.String(csvFieldDelimiterPrompt(*cfg.CSV.FieldDelimiter)),
				},
			}, nil
		}
	}
	return nil, xerrors.New("choose a input format type from: [json, csv]")
}

func outputFormatTypePrompt(current string) string {
	label := fmt.Sprintf("Select Output Format Type [%s]", current)
	return selectPrompt(label, []string{
		config.S3SELECT_INPUT_FORMAT_CSV,
		config.S3SELECT_INPUT_FORMAT_JSON,
	})
}

func selectPrompt(label string, items []string) string {
	p := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, _ := p.Run()
	return result
}

func KeySet(m map[string]string) []string {
	keyset := make([]string, 0, len(m))
	for k := range m {
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
