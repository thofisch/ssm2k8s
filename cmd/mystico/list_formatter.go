package main

import (
	"fmt"

	"github.com/thofisch/ssm2k8s/domain"
)

func NewListFormatter(secrets domain.ApplicationSecrets, verbose bool) func(args ...interface{}) {
	appNameMaxWidth := len("PATH")
	secretNameMaxWidth := len("SECRET")

	for secretName, secret := range secrets {
		length := len(secret.Path)
		if length > appNameMaxWidth {
			appNameMaxWidth = length
		}

		length = len(secretName)
		if length > secretNameMaxWidth {
			secretNameMaxWidth = length
		}

		if verbose {
			for name := range secret.Data {
				length := len(secret.Path + "/" + name)
				if length > appNameMaxWidth {
					appNameMaxWidth = length
				}
			}
		}
	}

	var format string

	if verbose {
		// PATH  SECRET  VERSION  LAST MODIFIED  VALUE
		format = fmt.Sprintf("%%-%ds  %%-%ds %%-7s  %%-20s  %%s\n", appNameMaxWidth, secretNameMaxWidth)
	} else {
		// PATH  SECRET  KEYS  HASH  LAST MODIFIED
		format = fmt.Sprintf("%%-%ds  %%-%ds  %%4s  %%-7s  %%-20s\n", appNameMaxWidth, secretNameMaxWidth)
	}

	return func(args ...interface{}) {
		fmt.Printf(format, args...)
	}
}
