package cmds

import (
	"fmt"
	"os"

	"github.com/ssoor/kuberes/cmds/build"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	colorable "github.com/mattn/go-colorable"

	"sigs.k8s.io/kustomize/pkg/factory"
	"sigs.k8s.io/kustomize/pkg/fs"
)

var (
	log = logrus.WithFields(logrus.Fields{"package": "cmds"})

	// RootCommand is the root of the command tree.
	RootCommand *cobra.Command
)

const (
	commandLongDesc = `
kustomize manages declarative configuration of Kubernetes.

See https://sigs.k8s.io/kustomize
`
)

// NewCommandRun is
func NewCommandRun(conf *viper.Viper, run func(cmd *cobra.Command, conf *viper.Viper, args []string)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		run(cmd, conf, args)
	}
}

// New returns an initialized command tree.
func New(f *factory.KustFactory) *cobra.Command {
	fsys := fs.MakeRealFS()
	stdOut := os.Stdout

	// Main dlv root command.
	RootCommand = &cobra.Command{
		Use:   "kuberes",
		Long:  commandLongDesc,
		Short: "kustomize manages declarative configuration of Kubernetes",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(colorable.NewColorableStdout())
			logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})

			switch cmd.Flag("log").Value.String() {
			case "debug":
				logrus.SetLevel(logrus.DebugLevel)
			case "info":
				logrus.SetLevel(logrus.InfoLevel)
			case "warn":
				logrus.SetLevel(logrus.WarnLevel)
			case "error":
				logrus.SetLevel(logrus.ErrorLevel)
			case "fatal":
				logrus.SetLevel(logrus.FatalLevel)
			case "panic":
				logrus.SetLevel(logrus.PanicLevel)
			default:
				return fmt.Errorf("unrecognized log output level")
			}

			return nil
		},
	}

	RootCommand.PersistentFlags().String("log", "info", "setting log output level.")

	RootCommand.AddCommand(
		build.NewCommand(stdOut, fsys, f.ResmapF, f.TransformerF),
	)

	RootCommand.DisableAutoGenTag = true

	return RootCommand
}
