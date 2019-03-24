// +build integration

package aws_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/thofisch/ssm2k8s/aws"
)

func TestGetParameters(t *testing.T) {
	stub := aws.NewParameterStore(endpoints.EuCentral1RegionID)

	parameters, _ := stub.GetParameters("/p-project/")

	result := NewResult(parameters)

	for k, v := range result.Secrets {
		fmt.Printf("Creating secret '\x1b[33m%s\x1b[0m' in namespace '\x1b[33m%s\x1b[0m' with the following values:\n", k, result.Name)
		alignKeyValue(v)
		fmt.Println()
	}

	for k, v := range result.Configs {
		fmt.Printf("Creating config-map '\x1b[33m%s\x1b[0m' in namespace '\x1b[33m%s\x1b[0m' with the following values:\n", k, result.Name)
		alignKeyValue(v)
		fmt.Println()
	}
}
