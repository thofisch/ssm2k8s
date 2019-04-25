package main

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"

	"github.com/thofisch/ssm2k8s/internal/logging"
	"gopkg.in/alecthomas/kingpin.v2"
)


type DeleteCommandOptions struct {
	Capability  string
	Environment string
	Application string
	Secrets     []string
}

func NewDeleteCommand(cmd *kingpin.CmdClause) *DeleteCommandOptions {
	opt := &DeleteCommandOptions{
		Secrets: []string{},
	}
	cmd.Flag("environment", "Environment").Short('e').Default("prod").StringVar(&opt.Environment)

	cmd.Arg("capability", "Nickname for user.").Required().StringVar(&opt.Capability)
	cmd.Arg("application", "Name of application.").Required().StringVar(&opt.Application)
	cmd.Arg("secrets", "listCmd of secrets").Required().StringsVar(&opt.Secrets)

	return opt
}

func ExecuteDelete(logger logging.Logger, options *DeleteCommandOptions) {
	parameterStore, err := aws.NewParameterStore(logger, "eu-central-1")
	if err != nil {
		panic(err)
	}

	for _, key := range options.Secrets {
		name := fmt.Sprintf("/%s/%s/%s/%s", options.Capability, options.Environment, options.Application, key)

		fmt.Printf("Deleting \033[33m%s\033[0m\n", name)
		err := parameterStore.DeleteApplicationSecret(options.Capability, options.Environment, options.Application, key)
		if err != nil {
			panic(err)
		}
	}
}
