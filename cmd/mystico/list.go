package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/internal/logging"
)

func ExecuteList(logger logging.Logger) {
	parameterStore, err := aws.NewParameterStore(logger, "eu-central-1")
	if err != nil {
		panic(err)
	}

	secrets, err := parameterStore.GetApplicationSecrets()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d secrets in %q\n", len(secrets), "AWS SSM Parameter Store")
	fmt.Printf("\n")

	appNameMaxWidth := 0
	keyNameMaxWidth := 0

	var secretNames []string

	for name, secret := range secrets {
		secretNames = append(secretNames, name)
		length := len(name)
		if length > appNameMaxWidth {
			appNameMaxWidth = length
		}
		for name := range secret.Data {
			length := len(name)
			if length > keyNameMaxWidth {
				keyNameMaxWidth = length
			}
		}
	}
	sf := fmt.Sprintf("%%-%ds          %%s     %%s", appNameMaxWidth)
	af := fmt.Sprintf("\033[34m%%-%ds\033[0m          %%s     %%s", appNameMaxWidth)
	hs := len(fmt.Sprintf(sf, "", "791efb9d8c0e74f81227afc39d4f24708f6aa8c3", "2019-03-21T12:43:27Z"))

	sort.Strings(secretNames)

	for _, appName := range secretNames {
		secret := secrets[appName]
		header := fmt.Sprintf(af,
			appName,
			secret.Hash,
			secret.LastModified.Format(time.RFC3339),
		)

		fmt.Printf(header + "\n")
		fmt.Printf(strings.Repeat("-", hs) + "\n")

		f := fmt.Sprintf("\033[33m%%-%ds\033[0m = \033[36m%%s\033[0m\n", keyNameMaxWidth)

		var keyNames []string

		for k := range secret.Data {
			keyNames = append(keyNames, k)
		}

		sort.Strings(keyNames)

		for _, k := range keyNames {
			v := secret.Data[k]
			fmt.Printf(f, k, v)
		}

		fmt.Printf("\n")
	}
}
