package aws

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"github.com/thofisch/ssm2k8s/internal/util"
)

type (
	ParameterStore interface {
		GetApplicationSecrets(capability string) (secrets domain.ApplicationSecrets, err error)
		PutApplicationSecret(capability string, environment string, application string, key string, value string, overwrite bool) error
		DeleteApplicationSecret(capability string, environment string, application string, key string) error
	}
	parameterStore struct {
		Log    logging.Logger
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

func NewParameterStore(logger logging.Logger, region string) (ParameterStore, error) {
	client, err := NewSsmClient(logger, NewSsmConfig(region))
	if err != nil {
		return nil, err
	}
	return NewParameterStoreWithClient(logger, client), nil
}

func NewParameterStoreWithClient(logger logging.Logger, client SsmClient) ParameterStore {
	return &parameterStore{
		Log:    logger,
		Client: client,
	}
}

func (ps *parameterStore) GetApplicationSecrets(capability string) (secrets domain.ApplicationSecrets, err error) {
	path := ensurePathPrefix(capability)

	ps.Log.Infof("Getting AWS SSM Parameters from Namespace %q", path)
	ssmParameters, err := ps.Client.GetParametersByPath(path)
	if err != nil {
		ps.Log.Errorf("[ERROR] %s\n", err)
		return nil, err
	}

	ps.Log.Debugf("Found %d parameters", len(ssmParameters))

	parameters := ps.filterParameters(ssmParameters)

	ps.Log.Debugf("Found %d parameters matching pattern /cap/env/app/key", len(ssmParameters))

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

func (ps *parameterStore) filterParameters(ssmParameters []*ssm.Parameter) []parameter {
	filteredParameters := make([]parameter, 0, len(ssmParameters))

	for _, p := range ssmParameters {
		name, err := parseParameterName(*p.Name)
		if err != nil {
			ps.Log.Debugf("Skipping incompatible parameter name %q", *p.Name)
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
		secretName := getSecretName(p.Name)
		m[secretName] = append(m[secretName], p)
	}

	return m
}

func getSecretName(pn parameterName) string {
	return fmt.Sprintf("%s-%s-secret", pn.Environment, pn.Application)
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

func (ps *parameterStore) PutApplicationSecret(capability string, environment string, application string, key string, value string, overwrite bool) error {
	name := fmt.Sprintf("/%s/%s/%s/%s", capability, environment, application, key)

	return ps.Client.PutParameter(name, value, overwrite)
}


func (ps *parameterStore) DeleteApplicationSecret(capability string, environment string, application string, key string) error {
	name := fmt.Sprintf("/%s/%s/%s/%s", capability, environment, application, key)

	return ps.Client.DeleteParameter(name)
}
