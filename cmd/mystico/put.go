package main

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

type PutCommandOptions struct {
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

	return opt
}

func ExecutePut(logger logging.Logger, options *PutCommandOptions) {
	parameterStore, err := aws.NewParameterStore(logger, *globalRegion)
	if err != nil {
		panic(err)
	}

	for key, value := range options.Secrets {
		application := strings.TrimLeft(options.Application, "/")
		fmt.Printf("Putting \033[33m/%s/%s\033[0m\n", application, key)

		err := parameterStore.PutApplicationSecret(application, key, value, options.Overwrite)
		if err != nil {
			panic(err)
		}
	}
}
