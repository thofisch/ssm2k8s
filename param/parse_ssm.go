package param

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"regexp"
	"time"
)

type ParameterName struct {
	Capability  string
	Environment string
	Application string
	Key         string
}

type ParameterInfo struct {
	Name         ParameterName
	Type         string
	Value        string
	LastModified time.Time
	Version      int64
}

type ParameterStore interface {
	GetParameters(path string) []ParameterInfo
}

type awsSystemManagerParameterStore struct {
	Region    string
	Recursive bool
	Decrypt   bool
}

func NewParameterStore(region string) ParameterStore {
	return &awsSystemManagerParameterStore{
		Region:    region,
		Recursive: true,
		Decrypt:   true,
	}
}

func (ps *awsSystemManagerParameterStore) GetParameters(path string) []ParameterInfo {
	fmt.Printf("Calling private\n")

	parametersByPath, err := ps.getParametersByPath(path)
	if err != nil {
		panic(err)
	}
	len := len(parametersByPath)
	var parameters = make([]ParameterInfo, len)

	for i, p := range parametersByPath {
		parameters[i], _ = toParameterInfo(p)
	}

	return parameters
}

func (ps *awsSystemManagerParameterStore) getParametersByPath(path string) (p []*ssm.Parameter, err error) {
	fmt.Printf("%#v\n", *ps)

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

func toParameterInfo(p *ssm.Parameter) (pi ParameterInfo, err error) {
	name, err := parseParameterName(*p.Name)
	if err != nil {
		return pi, err
	}

	pi.Name = name
	pi.Type = *p.Type
	pi.Value = *p.Value
	pi.LastModified = *p.LastModifiedDate
	pi.Version = *p.Version

	return pi, nil
}

var parameterNamePattern = regexp.MustCompile("^/(?P<cap>[^/]+)/(?P<env>[^/]+)/(?P<app>[^/]+)/(?P<key>[^/]+)$")

func parseParameterName(name string) (pn ParameterName, err error) {
	if !parameterNamePattern.MatchString(name) {
		return pn, errors.New(fmt.Sprintf("'%s' not of expected format: /cap/env/app/key", name))
	}

	groups := findNamedGroups(parameterNamePattern, name)

	pn.Capability = groups["cap"]
	pn.Environment = groups["env"]
	pn.Application = groups["app"]
	pn.Key = groups["key"]
	return pn, nil
}
