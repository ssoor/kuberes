package build

import (
	"errors"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"sigs.k8s.io/kustomize/pkg/constants"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/ifc/transformer"
	"sigs.k8s.io/kustomize/pkg/resmap"

	"github.com/ssoor/kuberes/pkg/loader"
	"github.com/ssoor/kuberes/pkg/target"
)

var (
	log = logrus.WithFields(logrus.Fields{"package": "onlyone"})
)

var examples = `
Use the file somedir/kustomization.yaml to generate a set of api resources:
    build somedir

Use a url pointing to a remote directory/kustomization.yaml to generate a set of api resources:
    build url
The url should follow hashicorp/go-getter URL format described in
https://github.com/hashicorp/go-getter#url-format

url examples:
  sigs.k8s.io/kustomize//examples/multibases?ref=v1.0.6
  github.com/Liujingfang1/mysql
  github.com/Liujingfang1/kustomize//examples/helloWorld?ref=repoUrl2
`

// NewCommand creates a new build command.
func NewCommand(out io.Writer, fs fs.FileSystem, rf *resmap.Factory, ptf transformer.Factory) *cobra.Command {
	ops := NewCommandOptions()

	cmd := &cobra.Command{
		Use:          "build [path]",
		Short:        "Print current configuration per contents of " + constants.KustomizationFileName,
		Example:      examples,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := ops.Validate(args)
			if err != nil {
				return err
			}
			return ops.Run(out, fs, rf, ptf)
		},
	}

	ops.SetFlags(cmd.Flags())

	return cmd
}

// CommandOptions contain the options for running a build
type CommandOptions struct {
	kustomizationPath string
	outputPath        string
}

// NewCommandOptions creates a BuildOptions object
func NewCommandOptions() *CommandOptions {
	return &CommandOptions{
		outputPath:        "",
		kustomizationPath: "./",
	}
}

// SetFlags set flags to command.
func (o *CommandOptions) SetFlags(flag *pflag.FlagSet) error {
	flag.StringVarP(&o.outputPath, "output", "o", "", "If specified, write the build output to this path.")

	return nil
}

// Validate validates command.
func (o *CommandOptions) Validate(args []string) error {
	if len(args) > 1 {
		return errors.New("specify one path to " + constants.KustomizationFileName)
	}
	if len(args) == 1 {
		o.kustomizationPath = args[0]
	}

	return nil
}

// Run runs command.
func (o *CommandOptions) Run(out io.Writer, fSys fs.FileSystem, rf *resmap.Factory, ptf transformer.Factory) error {
	ldr, err := loader.NewLoader(o.kustomizationPath, fSys)
	if err != nil {
		return err
	}
	defer ldr.Close()

	buildTarget, err := target.NewTarget(ldr)
	if err != nil {
		return err
	}

	if err := buildTarget.Load(); err != nil {
		return err
	}

	resourceMap, err := buildTarget.Make()
	if err != nil {
		return err
	}

	// Output the objects.
	res, err := resourceMap.Bytes()
	if err != nil {
		return err
	}

	if o.outputPath != "" {
		return fSys.WriteFile(o.outputPath, res)
	}
	_, err = out.Write(res)
	return err
}
