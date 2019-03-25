package aws

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/thofisch/ssm2k8s/internal/util"
)

type (
	ParameterStore interface {
		GetParameters(path string) ([]*parameter, error)
	}

	parameterStore struct {
		Client SsmClient
	}

	parameter struct {
		Name         *parameterName
		Value        ParameterValue
		LastModified time.Time
		Version      int64
	}

	parameterName struct {
		Capability  string
		Environment string
		Application string
		Key         string
	}
)

func NewParameterStore(region string) ParameterStore {
	return NewParameterStoreWithClient(NewSsmClient(NewSsmConfig(region)))
}

func NewParameterStoreWithClient(client SsmClient) ParameterStore {
	return &parameterStore{Client: client}
}

func (ps *parameterStore) GetParameters(path string) ([]*parameter, error) {
	samParameters, err := ps.Client.GetParametersByPath(path)
	if err != nil {
		return nil, err
	}

	parameters, err := mapParameters(samParameters)
	if err != nil {
		return nil, err
	}

	return parameters, nil
}

func mapParameters(ssmParameters []*ssm.Parameter) ([]*parameter, error) {
	var parameters = make([]*parameter, len(ssmParameters))

	for i, ssmParameter := range ssmParameters {
		parameter, err := mapParameter(ssmParameter)
		if err != nil {
			return nil, err
		}
		parameters[i] = parameter
	}

	return parameters, nil
}

func mapParameter(ssmParameter *ssm.Parameter) (*parameter, error) {
	parameterName, err := parseParameterName(*ssmParameter.Name)
	if err != nil {
		return nil, err
	}

	isSecret := isSecret(*ssmParameter.Type)
	parameterValue := NewParameterValue(*ssmParameter.Value, isSecret)

	return &parameter{
		Name:         parameterName,
		Value:        parameterValue,
		LastModified: *ssmParameter.LastModifiedDate,
		Version:      *ssmParameter.Version,
	}, nil
}

func isSecret(typeString string) bool {
	return ssm.ParameterTypeSecureString == typeString
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
