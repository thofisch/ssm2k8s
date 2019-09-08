package main

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"

	"github.com/thofisch/ssm2k8s/internal/logging"
	"gopkg.in/alecthomas/kingpin.v2"
)


type DeleteCommandOptions struct {
	Application string
	Environment string
	Secrets     []string
}

func NewDeleteCommand(cmd *kingpin.CmdClause) *DeleteCommandOptions {
	opt := &DeleteCommandOptions{
		Secrets: []string{},
	}

	cmd.Arg("application", "Name of application.").Required().StringVar(&opt.Application)
	cmd.Arg("secrets", "listCmd of secrets").Required().StringsVar(&opt.Secrets)

	cmd.Flag("environment", "Environment").Short('e').Default("prod").StringVar(&opt.Environment)

	return opt
}

func ExecuteDelete(logger logging.Logger, options *DeleteCommandOptions) {
	parameterStore, err := aws.NewParameterStore(logger, *globalRegion)
	if err != nil {
		panic(err)
	}

	for _, key := range options.Secrets {
		fmt.Printf("Deleting \033[33m/%s/%s/%s\033[0m\n", options.Application, options.Environment, key)

		err := parameterStore.DeleteApplicationSecret(options.Application, options.Environment, key)
		if err != nil {
			panic(err)
		}
	}
}
