package main

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"
	"strings"

	"github.com/thofisch/ssm2k8s/internal/logging"
	"gopkg.in/alecthomas/kingpin.v2"
)

type DeleteCommandOptions struct {
	Application string
	Secrets     []string
	Force       bool
}

func NewDeleteCommand(cmd *kingpin.CmdClause) *DeleteCommandOptions {
	opt := &DeleteCommandOptions{
		Secrets: []string{},
	}

	cmd.Arg("application", "Name of application.").Required().StringVar(&opt.Application)
	cmd.Arg("secrets", "listCmd of secrets").StringsVar(&opt.Secrets)

	cmd.Flag("force", "Force delete all secrets for application").Short('f').BoolVar(&opt.Force)

	return opt
}

func ExecuteDelete(logger logging.Logger, options *DeleteCommandOptions) {
	parameterStore, err := aws.NewParameterStore(logger, *globalRegion)
	if err != nil {
		panic(err)
	}

	application := strings.TrimLeft(options.Application, "/")

	if len(options.Secrets) == 0 {
		if !options.Force {
			panic(fmt.Errorf("requires the force flags when deleting all"))
		}

		secrets, err := parameterStore.GetApplicationSecrets(options.Application)
		if err != nil {
			panic(err)
		}

		for _, secret := range secrets {
			for key := range secret.Data {
				deleteParameter(parameterStore, application, key)
			}
		}
	} else {
		for _, key := range options.Secrets {
			deleteParameter(parameterStore, application, key)
		}
	}
}

func deleteParameter(parameterStore aws.ParameterStore, application string, key string) {
	fmt.Printf("Deleting \033[33m/%s/%s\033[0m\n", application, key)

	err := parameterStore.DeleteApplicationSecret(application, key)

	if err != nil {
		panic(err)
	}
}
