package aws

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type (
	ParameterStore interface {
		GetParameters(path string) ([]*Parameter, error)
	}

	Parameter struct {
		Name         *ParameterName
		Value        ParameterValue
		LastModified time.Time
		Version      int64
	}

	ParameterName struct {
		Capability  string
		Environment string
		Application string
		Key         string
	}

	awsParameterStore struct {
		Region    string
		Recursive bool
		Decrypt   bool
	}
)

func NewParameterStore(region string) ParameterStore {
	return &awsParameterStore{
		Region:    region,
		Recursive: true,
		Decrypt:   true,
	}
}

func (ps *awsParameterStore) GetParameters(path string) ([]*Parameter, error) {
	parameters, err := ps.getParametersByPath(path)
	if err != nil {
		return nil, err
	}

	return mapParameters(parameters)
}

func (ps *awsParameterStore) getParametersByPath(path string) ([]*ssm.Parameter, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(ps.Region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	client := ssm.New(sess, aws.NewConfig().WithRegion(ps.Region))

	output, err := client.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(ps.Recursive),
		WithDecryption: aws.Bool(ps.Decrypt),
	})
	if err != nil {
		return nil, err
	}

	return output.Parameters, nil
}

func (pn *ParameterName) String() string {
	return fmt.Sprintf("/%s/%s/%s/%s", pn.Capability, pn.Environment, pn.Application, pn.Key)
}

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
