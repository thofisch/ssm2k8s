package ssm2k8s

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/internal/util"
	"regexp"
)

func UpdateSecrets() {

	parameters, _ := aws.NewParameterStore("eu-central-1").GetParameters("/p-project/")

	for _, p := range parameters {
		fmt.Printf("%v\n", *p)
	}

}

type parameterName struct {
	Capability  string
	Environment string
	Application string
	Key         string
}

var parameterNamePattern = regexp.MustCompile("^/(?P<cap>[^/]+)/(?P<env>[^/]+)/(?P<app>[^/]+)/(?P<key>[^/]+)$")

func parseParameterName(name string) (*parameterName, error) {
	if !parameterNamePattern.MatchString(name) {
		return nil, fmt.Errorf("name '%s' is not of the expected format: /cap/env/app/key", name)
	}

	groups := util.FindNamedGroups(parameterNamePattern, name)

	return &parameterName{
		Capability:  groups["cap"],
		Environment: groups["env"],
		Application: groups["app"],
		Key:         groups["key"],
	}, nil
}

func (pn *parameterName) String() string {
	return fmt.Sprintf("/%s/%s/%s/%s", pn.Capability, pn.Environment, pn.Application, pn.Key)
}
