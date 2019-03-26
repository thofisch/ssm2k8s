package ssm2k8s

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	. "github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/internal/util"
)

func CreateParameters() {
	accountId, err := GetAccountId()
	if err != nil {
		panic(err)
	}

	path := "/p-project/"
	fmt.Printf("Creating AWS SystemManager Parameters in \033[33m%s\033[0m from account \033[33m%s\033[0m \n", path, accountId)

	newParameters := []struct {
		environment string
		application string
		key         string
		value       string
		secret      bool
	}{
		{"prod", "default", "kafka_bootstrap_servers", "pkc-l9pve.eu-west-1.aws.confluent.cloud:9092", true},
		{"prod", "default", "kafka_sasl_username", "2QIFZBFI32KT4I5Q", true},
		{"prod", "default", "kafka_sasl_password", "1gOjGYMZrDqHQe3mlYpytZP5ScnQGWkJnwUimIX8KkCK+4GHfDOOl+DnBjOMrFre", true},
	}

	for _, v := range newParameters {
		name := ParameterName{"p-project", v.environment, v.application, v.key}

		fmt.Printf("%s\n", &name)
	}

}

func GenerateSecretManifests() {
	accountId, err := GetAccountId()
	if err != nil {
		panic(err)
	}

	path := "/p-project/"
	fmt.Printf("Pulling AWS SystemManager Parameters in \033[33m%s\033[0m from account \033[33m%s\033[0m \n", path, accountId)

	parameters, _ := NewParameterStore("eu-central-1").GetParameters(path)
	parameterCount := len(parameters)
	if parameterCount == 0 {
		fmt.Printf(" No parameters found.")
		return
	}

	fmt.Printf(" Found %d parameters.\n\n", parameterCount)

	ps := findCompatibleParameters(parameters)
	result := NewResult(ps)

	namespace := result.Name

	for app, v := range result.Secrets {
		fmt.Printf(" Creating secret \033[33m%s\033[0m in namespace \033[33m%s\033[0m with:\n", app, namespace)

		for key, value := range v {
			fmt.Printf("  \033[33m%s\033[0m = \033[33m%s\033[0m\n", key, value)
		}

		fmt.Fprint(os.Stderr, printSecret(app, v))
		fmt.Fprintln(os.Stderr, "---")
	}

	for app, v := range result.Configs {
		fmt.Printf(" Creating secret \033[33m%s\033[0m in namespace \033[33m%s\033[0m with:\n", app, namespace)

		for key, value := range v {
			fmt.Printf("  \033[33m%s\033[0m = \033[33m%s\033[0m\n", key, value)
		}

		fmt.Fprint(os.Stderr, printSecret(app, v))
		fmt.Fprintln(os.Stderr, "---")
	}
}

func printSecret(name string, values ParameterValueMap) string {
	var sb strings.Builder

	_, err := fmt.Fprintf(&sb, `apiVersion: v1
kind: Secret
metadata:
  name: %s
type: Opaque
stringData:
`, name)
	if err != nil {
		panic(err)
	}

	for k, v := range values {
		_, err := fmt.Fprintf(&sb, "  %s: %s\n", k, v)
		if err != nil {
			panic(err)
		}
	}

	fmt.Fprintln(&sb)
	return sb.String()
}

func findCompatibleParameters(parameters []*Parameter) []*P {
	out := make([]*P, 0, len(parameters))

	for _, p := range parameters {
		name, err := parseParameterName(p.Name)
		if err != nil {
			fmt.Printf("\032mIncompatible parameter name '%s', skipping...\033{0m\n", p.Name)
			continue
		}

		out = append(out, &P{
			Name:     name,
			Value:    p.Value.GetValue(),
			IsSecret: p.Value.IsSecret(),
		})
	}

	return out
}

type ParameterName struct {
	Capability  string
	Environment string
	Application string
	Key         string
}

type P struct {
	Name     *ParameterName
	Value    string
	IsSecret bool
}

var parameterNamePattern = regexp.MustCompile("^/(?P<cap>[^/]+)/(?P<env>[^/]+)/(?P<app>[^/]+)/(?P<key>[^/]+)$")

func parseParameterName(name string) (*ParameterName, error) {
	if !parameterNamePattern.MatchString(name) {
		return nil, fmt.Errorf("name '%s' is not of the expected format: /cap/env/app/key", name)
	}

	groups := util.FindNamedGroups(parameterNamePattern, name)

	return &ParameterName{
		Capability:  groups["cap"],
		Environment: groups["env"],
		Application: groups["app"],
		Key:         groups["key"],
	}, nil
}

func (pn *ParameterName) String() string {
	return fmt.Sprintf("/%s/%s/%s/%s", pn.Capability, pn.Environment, pn.Application, pn.Key)
}

type ParameterValueMap map[string]string

type N map[string]ParameterValueMap

type Result struct {
	Name    string
	Secrets N
	Configs N
}

func NewResult(parameters []*P) Result {
	result := Result{Name: parameters[0].Name.Capability}
	result.Secrets = make(N)
	result.Configs = make(N)
	secrets := filterSecrets(parameters, true)
	configs := filterSecrets(parameters, false)
	{
		applications := applicationMap(secrets)
		for a, pi := range applications {
			keyValue := keyValueMap(pi)
			result.Secrets[a] = keyValue
		}
	}
	{
		applications := applicationMap(configs)
		for a, pi := range applications {
			keyValue := keyValueMap(pi)
			result.Configs[a] = keyValue
		}
	}
	return result
}

func filterSecrets(pi []*P, b bool) []*P {
	var a []*P

	for _, p := range pi {
		if p.IsSecret == b {
			a = append(a, p)
		}
	}

	return a
}

func applicationMap(pi []*P) map[string][]*P {
	m := make(map[string][]*P)

	for _, p := range pi {
		m[p.Name.Application] = append(m[p.Name.Application], p)
	}

	return m
}

func keyValueMap(pi []*P) ParameterValueMap {
	m := make(ParameterValueMap)

	for _, p := range pi {
		m[p.Name.Key] = p.Value
	}

	return m
}
