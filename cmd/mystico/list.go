package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ListCommandOptions struct {
	Application string
	Verbose     bool
	Decode      bool
}

func NewListCommand(cmd *kingpin.CmdClause) *ListCommandOptions {
	opt := &ListCommandOptions{}

	//cmd.Arg("application", "Name of application.").HintAction(applicationHint).StringVar(&opt.Application)

	cmd.Flag("verbose", "Print keys and values").Short('v').BoolVar(&opt.Verbose)
	cmd.Flag("decode", "Print decoded values").Short('d').BoolVar(&opt.Decode)

	return opt
}

//func applicationHint() []string {
//	return []string{"p-project", "control-tower"}
//}

func ExecuteList(logger logging.Logger, opt *ListCommandOptions) {
	parameterStore, err := aws.NewParameterStore(logger, *globalRegion)
	if err != nil {
		panic(err)
	}

	secrets, err := parameterStore.GetApplicationSecrets()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d secrets in %q\n", len(secrets), "AWS SSM Parameter Store")
	fmt.Printf("\n")

	format := NewListFormatter(secrets, opt.Verbose)

	if opt.Verbose {
		format("SECRET", "KEYS", "HASH", "LAST MODIFIED", "KEY", "VALUE")
	} else {
		format("SECRET", "KEYS", "HASH", "LAST MODIFIED")
	}

	for _, secretName := range sortApplicationSecrets(secrets) {
		secret := secrets[secretName]

		if opt.Verbose {
			for i, key := range sortSecretData(secret.Data) {

				var value string

				if opt.Decode {
					value = secret.Data[key]
				} else {
					value = "***"
				}

				if i == 0 {
					format(secret.Path, strconv.Itoa(len(secret.Data)), secret.Hash[0:7], secret.LastModified.Format(time.RFC3339), key, value)
				} else {
					format("", "", "", "", key, value)
				}
			}

		} else {
			format(secret.Path, strconv.Itoa(len(secret.Data)), secret.Hash[0:7], secret.LastModified.Format(time.RFC3339))
		}
	}
}

func sortApplicationSecrets(secrets domain.ApplicationSecrets) []string {
	var secretNames []string
	for name := range secrets {
		secretNames = append(secretNames, name)
	}
	sort.Strings(secretNames)
	return secretNames
}

func sortSecretData(data domain.SecretData) []string {
	var keyNames []string
	for k := range data {
		keyNames = append(keyNames, k)
	}
	sort.Strings(keyNames)
	return keyNames
}
