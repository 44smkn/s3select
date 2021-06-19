package root

import (
	"github.com/44smkn/s3select/pkg/cli"
	configureCmd "github.com/44smkn/s3select/pkg/cli/configure"
	queryCmd "github.com/44smkn/s3select/pkg/cli/query"
	versionCmd "github.com/44smkn/s3select/pkg/cli/version"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdRoot(f *cli.Factory, version, buildDate string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "s3select <command> [flags]",
		Short: "S3Select CLI implemented with golang",
		Long:  `Work seamlessly with s3select from the command line.`,

		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
			$ s3select configure
			$ s3select query -b bucket -k path/to/prefix -p default -s 'SELECT * LIMIT 10'
		`),
	}

	cmd.SetOut(f.Out)
	cmd.SetErr(f.ErrOut)

	cmd.AddCommand(versionCmd.NewCmdVersion(version, buildDate))
	cmd.AddCommand(configureCmd.NewCmdCongigure(f))
	cmd.AddCommand(queryCmd.NewCmdQuery(f))

	cmd.PersistentFlags().Bool("help", false, "Show help for command")

	return cmd
}
