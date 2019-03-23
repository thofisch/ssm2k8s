package param_test

import (
	"fmt"
	. "github.com/thofisch/ssm2k8s/param"
	"testing"
	"time"
)

type Result struct {
	Name    string
	Secrets map[string]Secret
	Configs map[string]Config
}

type Secret struct {
	Data map[string]ParameterValue
}

type Config struct {
	Data map[string]ParameterValue
}

func Test_aName(t *testing.T) {

	stub := NewParameterStoreStub(
		aParameterInfo(
			withName(aParameterName(
				withApplication("foo"),
				withKey("pghost"))),
			withSecret("lala")),
		aParameterInfo(
			withName(aParameterName(
				withApplication("foo"),
				withKey("pguser"))),
			withSecret("lala")),
		aParameterInfo(
			withName(aParameterName(
				withApplication("foo"),
				withKey("pgpassword"))),
			withSecret("lala")),
		aParameterInfo(
			withName(aParameterName(
				withApplication("foo"),
				withKey("pgport"))),
			withValue("1433")),
		aParameterInfo(
			withName(aParameterName(withKey("kafka-brokers"))),
			withValue("broker1,broker2,broker3")),
	)

	parameters, _ := stub.GetParameters("/p-project/")

	result := NewResult(parameters)

	for k, v := range result.Secrets {
		fmt.Printf("Creating secret '\x1b[33m%s\x1b[0m' in namespace '\x1b[33m%s\x1b[0m' with the following values:\n", k, result.Name)
		alignKeyValue(v.Data)
		fmt.Println()
	}

	for k, v := range result.Configs {
		fmt.Printf("Creating config-map '\x1b[33m%s\x1b[0m' in namespace '\x1b[33m%s\x1b[0m' with the following values:\n", k, result.Name)
		alignKeyValue(v.Data)
		fmt.Println()
	}
}

func alignKeyValue(m map[string]ParameterValue) {
	var max = 0
	for k := range m {
		l := len(k)
		if l > max {
			max = l
		}
	}

	format := fmt.Sprintf("  \x1b[33m%%-%ds\x1b[0m = \x1b[34m%%s\x1b[0m\n", max)

	for k, v := range m {
		fmt.Printf(format, k, v)
	}
}

func NewResult(parameters []*ParameterInfo) Result {
	result := Result{Name: parameters[0].Name.Capability}
	result.Secrets = make(map[string]Secret)
	result.Configs = make(map[string]Config)
	secrets := fitlerSecrets(parameters, true)
	configs := fitlerSecrets(parameters, false)
	{
		applications := applicationMap(secrets)
		for a, pi := range applications {
			keyValue := keyValueMap(pi)
			result.Secrets[a] = Secret{Data: keyValue}
		}
	}
	{
		applications := applicationMap(configs)
		for a, pi := range applications {
			keyValue := keyValueMap(pi)
			result.Configs[a] = Config{Data: keyValue}
		}
	}
	return result
}

func applicationMap(pi []*ParameterInfo) map[string][]*ParameterInfo {
	m := make(map[string][]*ParameterInfo)

	for _, p := range pi {
		m[p.Name.Application] = append(m[p.Name.Application], p)
	}

	return m
}

func keyValueMap(pi []*ParameterInfo) map[string]ParameterValue {
	m := make(map[string]ParameterValue)

	for _, p := range pi {
		m[p.Name.Key] = p.Value
	}

	return m
}

func fitlerSecrets(pi []*ParameterInfo, b bool) []*ParameterInfo {
	var a []*ParameterInfo

	for _, p := range pi {
		if p.Value.IsSecret() == b {
			a = append(a, p)
		}
	}

	return a
}

type ParameterStoreStub struct {
	parameterInfoList []*ParameterInfo
}

func NewParameterStoreStub(parameters ...*ParameterInfo) ParameterStore {
	return &ParameterStoreStub{parameterInfoList: parameters}
}

func (ps *ParameterStoreStub) GetParameters(path string) ([]*ParameterInfo, error) {
	return ps.parameterInfoList, nil
}

func aParameterInfo(builders ...func(*ParameterInfo)) *ParameterInfo {
	pi := &ParameterInfo{
		Name:         aParameterName(),
		Value:        NewParameterValue("", false),
		Version:      0,
		LastModified: time.Time{},
	}

	for _, builder := range builders {
		builder(pi)
	}

	return pi
}

func withName(pn *ParameterName) func(*ParameterInfo) {
	return func(p *ParameterInfo) {
		p.Name = pn
	}
}

func withValue(v string) func(*ParameterInfo) {
	return func(p *ParameterInfo) {
		p.Value = NewParameterValue(v, false)
	}
}

func withSecret(v string) func(*ParameterInfo) {
	return func(p *ParameterInfo) {
		p.Value = NewParameterValue(v, true)
	}
}

func aParameterName(builders ...func(*ParameterName)) *ParameterName {
	pn := &ParameterName{
		Capability:  "selfservice",
		Environment: "prod",
		Application: "default",
		Key:         "",
	}

	for _, builder := range builders {
		builder(pn)
	}

	return pn
}

func withApplication(application string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Application = application
	}
}

func withKey(key string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Key = key
	}
}
