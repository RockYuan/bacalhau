package validate

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// NewCmd creates and returns a new validate command
func NewCmd() *cobra.Command {
	opts := &ValidateOptions{}

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a job using a JSON or YAML file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Validate(opts); err != nil {
				return err
			}
			cmd.Println("The jobspec is valid")
			return nil
		},
	}

	validateCmd.Flags().StringVarP(&opts.Filename, "filename", "f", "", "File containing the job to validate")
	validateCmd.MarkFlagRequired("filename")

	return validateCmd
}

type ValidateOptions struct {
	Filename string
}

func Validate(opts *ValidateOptions) error {
	cueCtx := cuecontext.New()

	buildInstances := load.Instances([]string{"job-schema.cue"}, &load.Config{
		Dir: "pkg/models",
	})

	if len(buildInstances) == 0 {
		return fmt.Errorf("no CUE instances found")
	}

	if buildInstances[0].Err != nil {
		return fmt.Errorf("error loading CUE instance: %w", buildInstances[0].Err)
	}

	instance := cueCtx.BuildInstance(buildInstances[0])
	if instance.Err() != nil {
		return fmt.Errorf("error building CUE instance: %w", instance.Err())
	}

	jobSchema := instance.LookupPath(cue.ParsePath("#Job"))
	if !jobSchema.Exists() {
		return fmt.Errorf("#Job schema not found in the CUE model")
	}

	data, err := os.ReadFile(opts.Filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", opts.Filename, err)
	}

	jsonBytes, err := yaml.YAMLToJSON(data)
	if err != nil {
		return fmt.Errorf("error converting YAML to JSON: %w", err)
	}

	value := cueCtx.CompileBytes(jsonBytes, cue.Filename(opts.Filename))
	unified := value.Unify(jobSchema)
	if err := unified.Validate(cue.Concrete(true), cue.Final()); err != nil {
		return err
	}

	return nil
}
