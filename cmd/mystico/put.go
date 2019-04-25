package main

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"

	"github.com/thofisch/ssm2k8s/internal/logging"
	"gopkg.in/alecthomas/kingpin.v2"
)

type PutCommandOptions struct {
	Capability  string
	Environment string
	Application string
	Secrets     map[string]string
	Overwrite   bool
}

func NewPutCommand(cmd *kingpin.CmdClause) *PutCommandOptions {
	opt := &PutCommandOptions{
		Secrets: map[string]string{},
	}
	cmd.Flag("overwrite", "Overwrite existing secrets").Short('o').BoolVar(&opt.Overwrite)
	cmd.Flag("environment", "Environment").Short('e').Default("prod").StringVar(&opt.Environment)

	cmd.Arg("capability", "Nickname for user.").Required().StringVar(&opt.Capability)
	cmd.Arg("application", "Name of application.").Required().StringVar(&opt.Application)
	cmd.Arg("secrets", "listCmd of secrets").Required().StringMapVar(&opt.Secrets)

	return opt
}

func ExecutePut(logger logging.Logger, options *PutCommandOptions) {
	parameterStore, err := aws.NewParameterStore(logger, "eu-central-1")
	if err != nil {
		panic(err)
	}

	for key, value := range options.Secrets {
		name := fmt.Sprintf("/%s/%s/%s/%s", options.Capability, options.Environment, options.Application, key)

		fmt.Printf("Putting \033[33m%s\033[0m\n", name)
		err := parameterStore.PutApplicationSecret(options.Capability, options.Environment, options.Application, key, value, options.Overwrite)
		if err != nil {
			panic(err)
		}
	}

}
