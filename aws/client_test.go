// +build integration

package aws_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/thofisch/ssm2k8s/aws"
)

func TestGetParameters(t *testing.T) {
	stub := aws.NewParameterStore("eu-central-1")

	parameters, _ := stub.GetParameters("/p-project/")

	for _, p := range parameters {
		fmt.Printf("%s %s\n", p.LastModified.Format(time.RFC3339), p.LastModified.Local().Format(time.RFC3339))
	}
}
