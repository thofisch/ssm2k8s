package aws

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/internal/util"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
)

type (
	ParameterStore interface {
		GetApplicationSecrets(capability string) (secrets domain.ApplicationSecrets, err error)
	}
	parameterStore struct {
		Client SsmClient
	}
	parameter struct {
		Name         parameterName
		Value        string
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

func NewParameterStore(region string) (ParameterStore, error) {
	client, err := NewSsmClient(NewSsmConfig(region))
	if err != nil {
		return nil, err
	}
	return NewParameterStoreWithClient(client), nil
}

func NewParameterStoreWithClient(client SsmClient) ParameterStore {
	return &parameterStore{Client: client}
}

func (ps *parameterStore) GetApplicationSecrets(capability string) (secrets domain.ApplicationSecrets, err error) {
	path := ensurePathPrefix(capability)
	ssmParameters, err := ps.Client.GetParametersByPath(path)
	if err != nil {
		return nil, err
	}

	//awsParameters := mapParameters(ssmParameters)

	//parameterCount := len(awsParameters)
	//if parameterCount == 0 {
	//	fmt.Printf(" No awsParameters found.")
	//	var parameters []*parameter = make([]*parameter, 0)
	//	return getApplicationSecrets(parameters), nil
	//}
	//fmt.Printf(" Found %d awsParameters.\n\n", parameterCount)

	parameters := filterParameters(ssmParameters)

	secrets = getApplicationSecrets(parameters)

	return
}

func ensurePathPrefix(s string) string {
	if strings.HasPrefix(s, "/") {
		return s
	} else {
		return "/" + s
	}
}

func filterParameters(ssmParameters []*ssm.Parameter) []parameter {
	filteredParameters := make([]parameter, 0, len(ssmParameters))

	for _, p := range ssmParameters {
		name, err := parseParameterName(*p.Name)
		if err != nil {
			fmt.Printf("\032mIncompatible parameter name '%s', skipping...\033{0m\n", *p.Name)
			continue
		}

		filteredParameters = append(filteredParameters, parameter{
			Name:         name,
			Value:        *p.Value,
			LastModified: *p.LastModifiedDate,
			Version:      *p.Version,
		})
	}

	return filteredParameters
}

var parameterNamePattern = regexp.MustCompile("^/(?P<cap>[^/]+)/(?P<env>[^/]+)/(?P<app>[^/]+)/(?P<key>[^/]+)$")

func parseParameterName(name string) (parameterName, error) {
	if !parameterNamePattern.MatchString(name) {
		return parameterName{}, fmt.Errorf("name '%s' is not of the expected format: /capability/environment/application/key", name)
	}

	groups := util.FindNamedGroups(parameterNamePattern, name)

	return parameterName{
		Capability:  groups["cap"],
		Environment: groups["env"],
		Application: groups["app"],
		Key:         groups["key"],
	}, nil
}

func getApplicationSecrets(parameters []parameter) domain.ApplicationSecrets {
	secrets := make(domain.ApplicationSecrets)

	applications := mapApplications(parameters)

	for appName, appParameters := range applications {
		data := mapData(appParameters)
		secrets[appName] = domain.ApplicationSecret{
			LastModified: util.FindNewest(getDates(appParameters)),
			Hash:         util.HashKeyValuePairs(getKeyValuePairs(appParameters)),
			Data:         data,
		}
	}

	return secrets
}

func mapApplications(parameters []parameter) map[string][]parameter {
	m := make(map[string][]parameter)

	for _, p := range parameters {
		m[p.Name.Application] = append(m[p.Name.Application], p)
	}

	return m
}

func mapData(parameters []parameter) domain.SecretData {
	secretData := make(domain.SecretData)

	for _, p := range parameters {
		secretData[p.Name.Key] = p.Value
	}

	return secretData
}

func getDates(parameters []parameter) []time.Time {
	dates := make([]time.Time, 0, len(parameters))
	for _, v := range parameters {
		dates = append(dates, v.LastModified)
	}
	return dates
}

func getKeyValuePairs(parameters []parameter) map[string]string {
	kv := make(map[string]string, len(parameters))

	for _, p := range parameters {
		kv[p.Name.Key] = p.Value
	}

	return kv
}
