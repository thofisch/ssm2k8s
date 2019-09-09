package main

import (
	"fmt"

	"github.com/thofisch/ssm2k8s/domain"
)

func NewListFormatter(secrets domain.ApplicationSecrets, verbose bool) func (args ...interface{}) {
	appNameMaxWidth := 0
	keyNameMaxWidth := 0

	for _, secret := range secrets {
		length := len(secret.Path)
		if length > appNameMaxWidth {
			appNameMaxWidth = length
		}

		if verbose {
			for name := range secret.Data {
				length := len(name)
				if length > keyNameMaxWidth {
					keyNameMaxWidth = length
				}
			}
		}
	}

	var format string

	if verbose {
		format = fmt.Sprintf("%%-%ds  %%4s  %%-7s  %%-20s | %%-%ds %%s\n", appNameMaxWidth, keyNameMaxWidth)
	} else {
		format = fmt.Sprintf("%%-%ds  %%4s  %%-7s  %%-20s\n", appNameMaxWidth)
	}

	return func(args ...interface{}) {
		fmt.Printf(format, args...)
	}
}
