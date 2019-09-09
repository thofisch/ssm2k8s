package aws

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/internal/config"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"github.com/thofisch/ssm2k8s/internal/util"
)

type (
	ParameterStore interface {
		GetApplicationSecrets() (secrets domain.ApplicationSecrets, err error)
		PutApplicationSecret(application string, key string, value string, overwrite bool) error
		DeleteApplicationSecret(application string, key string) error
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
		Application string
		Paths       []string
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

func (ps *parameterStore) GetApplicationSecrets() (secrets domain.ApplicationSecrets, err error) {
	ps.Log.Info("Getting AWS SSM Parameters")

	ssmParameters, err := ps.Client.GetParametersByPath("/")
	if err != nil {
		ps.Log.Errorf("[ERROR] %s\n", err)
		return nil, err
	}

	parameters := ps.filterParameters(ssmParameters)

	ps.Log.Debugf("Found %d parameters matching pattern %s", len(parameters), expectedFormat)

	secrets = getApplicationSecrets(parameters)

	return
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

const expectedFormat = "/application/*/key"

var parameterNamePattern = regexp.MustCompile("^/(?P<app>[^/]+)(?:/?(?P<paths>.+))?/(?P<key>[^/]+)$")

func parseParameterName(name string) (parameterName, error) {
	if !parameterNamePattern.MatchString(name) {
		return parameterName{}, fmt.Errorf("name %q is not of the expected format %q", name, expectedFormat)
	}

	groups := util.FindNamedGroups(parameterNamePattern, name)
	application := groups["app"]

	if config.IsReservedWord(application) {
		return parameterName{}, fmt.Errorf("ignored %q as it uses the one of the reservered words: %q", name, config.ReservedWords)
	}

	pathParts, err := getPaths(groups["paths"], name)

	if err != nil {
		return parameterName{}, err
	}

	return parameterName{
		Application: application,
		Paths:       pathParts,
		Key:         groups["key"],
	}, nil
}

func getPaths(paths string, name string) ([]string, error) {
	if paths != "" {
		pathParts := strings.Split(paths, "/")

		for _, path := range pathParts {
			if len(path) <= 0 {
				return nil, fmt.Errorf("ignored %q as it has an illegal path", name)
			}
		}

		return pathParts, nil
	} else {
		return []string{}, nil
	}
}

func getApplicationSecrets(parameters []parameter) domain.ApplicationSecrets {
	secrets := make(domain.ApplicationSecrets)

	applications := mapApplications(parameters)

	for appName, appParameters := range applications {
		data := mapData(appParameters)
		secrets[appName] = domain.ApplicationSecret{
			Path:         getPath(appParameters[0].Name),
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
	return strings.Join(append([]string{pn.Application}, pn.Paths...), "-")
}

func getPath(pn parameterName) string {
	return fmt.Sprintf("/%s", strings.Join(append([]string{pn.Application}, pn.Paths...), "/"))
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

func (ps *parameterStore) PutApplicationSecret(application string, key string, value string, overwrite bool) error {
	name := fmt.Sprintf("/%s/%s", application, key)

	return ps.Client.PutParameter(name, value, overwrite)
}

func (ps *parameterStore) DeleteApplicationSecret(application string, key string) error {
	name := fmt.Sprintf("/%s/%s", application, key)

	return ps.Client.DeleteParameter(name)
}
