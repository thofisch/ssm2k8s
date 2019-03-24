package param

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ssm"
	"regexp"
)

func mapParameters(ssmParameters []*ssm.Parameter) ([]*Parameter, error) {
	var parameters = make([]*Parameter, len(ssmParameters))

	for i, ssmParameter := range ssmParameters {
		parameter, err := mapParameter(ssmParameter)
		if err != nil {
			return nil, err
		}
		parameters[i] = parameter
	}

	return parameters, nil
}

func mapParameter(ssmParameter *ssm.Parameter) (*Parameter, error) {
	parameterName, err := parseParameterName(*ssmParameter.Name)
	if err != nil {
		return nil, err
	}

	isSecret := isSecret(*ssmParameter.Type)
	parameterValue := NewParameterValue(*ssmParameter.Value, isSecret)

	return &Parameter{
		Name:         parameterName,
		Value:        parameterValue,
		LastModified: *ssmParameter.LastModifiedDate,
		Version:      *ssmParameter.Version,
	}, nil
}

var parameterNamePattern = regexp.MustCompile("^/(?P<cap>[^/]+)/(?P<env>[^/]+)/(?P<app>[^/]+)/(?P<key>[^/]+)$")

func parseParameterName(name string) (*ParameterName, error) {
	if !parameterNamePattern.MatchString(name) {
		return nil, fmt.Errorf("name '%s' is not of the expected format: /cap/env/app/key", name)
	}

	groups := findNamedGroups(parameterNamePattern, name)

	return &ParameterName{
		Capability:  groups["cap"],
		Environment: groups["env"],
		Application: groups["app"],
		Key:         groups["key"],
	}, nil
}

func isSecret(typeString string) bool {
	return ssm.ParameterTypeSecureString == typeString
}
