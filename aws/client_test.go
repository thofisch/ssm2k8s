// +build integration

package aws_test

import (
	"fmt"
	"testing"

	"github.com/thofisch/ssm2k8s/aws"
)

func TestGetParameters(t *testing.T) {
	stub := aws.NewSsmClient(aws.NewSsmConfig("eu-central-1"))

	parameters, _ := stub.GetParameters("/p-project/")

	for _, p := range parameters {
		fmt.Printf("%v\n", *p)
	}
}
