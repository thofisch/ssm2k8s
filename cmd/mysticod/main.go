package main

import (
	"github.com/thofisch/ssm2k8s"
	"github.com/thofisch/ssm2k8s/k8s"
)

func main() {

	ssm2k8s.UpdateSecrets()

	k8s.GetSecrets()

}

