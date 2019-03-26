package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
)

type (
	ParameterStore interface {
		GetParameters(path string) ([]*Parameter, error)
	}

	parameterStore struct {
		Client SsmClient
	}

	Parameter struct {
		Name         string
		Value        ParameterValue
		LastModified time.Time
		Version      int64
	}
)

func NewParameterStore(region string) ParameterStore {
	return NewParameterStoreWithClient(NewSsmClient(NewSsmConfig(region)))
}

func NewParameterStoreWithClient(client SsmClient) ParameterStore {
	return &parameterStore{Client: client}
}

func (ps *parameterStore) GetParameters(path string) ([]*Parameter, error) {
	samParameters, err := ps.Client.GetParametersByPath(path)
	if err != nil {
		return nil, err
	}

	parameters := mapParameters(samParameters)

	return parameters, nil
}

func mapParameters(ssmParameters []*ssm.Parameter) []*Parameter {
	var parameters = make([]*Parameter, len(ssmParameters))

	for i, ssmParameter := range ssmParameters {
		parameters[i] = mapParameter(ssmParameter)
	}

	return parameters
}

func mapParameter(ssmParameter *ssm.Parameter) *Parameter {
	isSecret := isSecret(ssmParameter.Type)
	parameterValue := NewParameterValue(*ssmParameter.Value, isSecret)

	return &Parameter{
		Name:         *ssmParameter.Name,
		Value:        parameterValue,
		LastModified: *ssmParameter.LastModifiedDate,
		Version:      *ssmParameter.Version,
	}
}

func isSecret(typeString *string) bool {
	return ssm.ParameterTypeSecureString == *typeString
}
