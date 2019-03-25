package ssm2k8s

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"
)

func UpdateSecrets() {

	parameters, _ := aws.NewParameterStore("eu-central-1").GetParameters("/p-project/")

	for _, p := range parameters {
		fmt.Printf("%v\n", *p)
	}

}
