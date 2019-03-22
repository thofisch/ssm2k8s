package param

import (
	"fmt"
	"regexp"
)

type ParameterName struct {
	Capability  string
	Environment string
	Application string
	Key         string
}

func (pn *ParameterName) String() string {
	return fmt.Sprintf("/%s/%s/%s/%s", pn.Capability, pn.Environment, pn.Application, pn.Key)
}

var parameterNamePattern = regexp.MustCompile("^/(?P<cap>[^/]+)/(?P<env>[^/]+)/(?P<app>[^/]+)/(?P<key>[^/]+)$")

func parseParameterName(name string) (*ParameterName, error) {
	if !parameterNamePattern.MatchString(name) {
		return nil, fmt.Errorf("name '%s' is not of expected format: /cap/env/app/key", name)
	}

	groups := findNamedGroups(parameterNamePattern, name)

	pn := &ParameterName{}
	pn.Capability = groups["cap"]
	pn.Environment = groups["env"]
	pn.Application = groups["app"]
	pn.Key = groups["key"]
	return pn, nil
}
