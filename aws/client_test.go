// +build integration

package aws_test

import (
	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/domain"
	"testing"
)

func TestGetParameters(t *testing.T) {
	stub,_ := aws.NewParameterStore("eu-central-1")

	secrets, _ := stub.GetApplicationSecrets("p-project")

	domain.PrintApplicationSecrets(secrets, "test")
}
