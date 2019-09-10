package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"github.com/thofisch/ssm2k8s/internal/util"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ListCommandOptions struct {
	Application string
	Verbose     bool
	Decode      bool
}

func NewListCommand(cmd *kingpin.CmdClause) *ListCommandOptions {
	opt := &ListCommandOptions{}

	cmd.Arg("application", "Name of application.") /*.HintAction(applicationHint)*/ .StringVar(&opt.Application)

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

	secrets, err := parameterStore.GetApplicationSecrets(opt.Application)
	if err != nil {
		panic(err)
	}

	numSecrets := len(secrets)
	fmt.Printf("Found %d secrets in %q\n", numSecrets, "AWS SSM Parameter Store")

	if numSecrets == 0 {
		return
	}

	fmt.Printf("\n")

	columnizer := util.NewColumnizer()

	if opt.Verbose {
		columnizer.Append("PATH", "SECRET", "VERSION", "LAST MODIFIED", "VALUE")
	} else {
		columnizer.Append("PATH", "SECRET", "KEYS", "HASH", "LAST MODIFIED")
	}

	for _, secretName := range sortApplicationSecrets(secrets) {
		secret := secrets[secretName]

		if opt.Verbose {
			for _, key := range sortSecretData(secret.Data) {
				dataSecret := secret.Data[key]
				path := secret.Path + "/" + key
				version := strconv.FormatInt(dataSecret.Version, 10)
				lastModified := dataSecret.LastModified.Format(time.RFC3339)
				value:= dataSecret.Value

				if !opt.Decode {
					value = "***"
				}

				columnizer.Append(path, secretName, version, lastModified, value)
			}
		} else {
			columnizer.Append(secret.Path, secretName, strconv.Itoa(len(secret.Data)), secret.Hash[0:7], secret.LastModified.Format(time.RFC3339))
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	columnizer.Print(writer)
	defer writer.Flush()
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
