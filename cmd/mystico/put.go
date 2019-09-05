package main

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"

	"github.com/thofisch/ssm2k8s/internal/logging"
	"gopkg.in/alecthomas/kingpin.v2"
)

type PutCommandOptions struct {
	Environment string
	Application string
	Secrets     map[string]string
	Overwrite   bool
}

func NewPutCommand(cmd *kingpin.CmdClause) *PutCommandOptions {
	opt := &PutCommandOptions{
		Secrets: map[string]string{},
	}

	cmd.Arg("application", "Name of application.").Required().StringVar(&opt.Application)
	cmd.Arg("secrets", "listCmd of secrets").Required().StringMapVar(&opt.Secrets)

	cmd.Flag("overwrite", "Overwrite existing secrets").Short('o').BoolVar(&opt.Overwrite)
	cmd.Flag("environment", "Environment").Short('e').HintOptions("foo", "bar", "baz").Default("prod").StringVar(&opt.Environment)

	return opt
}

func ExecutePut(logger logging.Logger, options *PutCommandOptions) {
	parameterStore, err := aws.NewParameterStore(logger, "eu-central-1")
	if err != nil {
		panic(err)
	}

	for key, value := range options.Secrets {
		fmt.Printf("Putting \033[33m/%s/%s/%s\033[0m\n", options.Application, options.Environment, key)

		err := parameterStore.PutApplicationSecret(options.Application, options.Environment, key, value, options.Overwrite)
		if err != nil {
			panic(err)
		}
	}

}
