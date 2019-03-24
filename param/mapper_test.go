package param

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func Test_mapParameters(t *testing.T) {
	value := "val"
	lastModified, _ := time.Parse("2006-01-02T15:04:05", "2019-01-01T00:00:00")
	version := int64(1)
	parameter := &ssm.Parameter{
		Name:             aws.String("/cap/env/app/key"),
		Value:            aws.String(value),
		Version:          aws.Int64(version),
		Type:             aws.String(ssm.ParameterTypeSecureString),
		LastModifiedDate: aws.Time(lastModified),
	}

	result, err := mapParameters([]*ssm.Parameter{parameter})

	assertOk(t, err)
	assertEqual(t, result[0], AParameter(
		WithName(AParameterName()),
		WithSecret(value),
		WithVersion(version),
		WithLastModified(lastModified),
	))
}

func Test_parseParameterName_valid_name(t *testing.T) {
	pn, err := parseParameterName("/a/b/c/d")

	assertOk(t, err)
	assertEqual(t, "a", pn.Capability)
	assertEqual(t, "b", pn.Environment)
	assertEqual(t, "c", pn.Application)
	assertEqual(t, "d", pn.Key)
}

func Test_parseParameterName_invalid_name(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "only slashes", input: "////"},
		{name: "missing slash prefix", input: "a/b/c/d/e"},
		{name: "slash_cap", input: "/a"},
		{name: "slash_cap_slash", input: "/a/"},
		{name: "slash_cap_slash_env", input: "/a/b"},
		{name: "slash_cap_slash_env_slash", input: "/a/b/"},
		{name: "slash_cap_slash_env_slash_app", input: "/a/b/c"},
		{name: "slash_cap_slash_env_slash_app_slash", input: "/a/b/c/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash", input: "/a/b/c/d/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash_extra", input: "/a/b/c/d/e"},
		{name: "slash_cap_slash_env_slash_slash_key", input: "/a/b//d"},
		{name: "slash_cap_slash_slash app_slash_key", input: "/a//c/d"},
		{name: "slash_slash_env_slash_app_slash_key", input: "//b/c/d"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseParameterName(test.input)

			assertNotOk(t, err)
		})
	}
}

func AParameter(builders ...func(*Parameter)) *Parameter {
	var pi = &Parameter{
		Name:         AParameterName(),
		Value:        NewParameterValue("", false),
		Version:      0,
		LastModified: time.Time{},
	}

	for _, builder := range builders {
		builder(pi)
	}

	return pi
}

func WithName(pn *ParameterName) func(*Parameter) {
	return func(p *Parameter) {
		p.Name = pn
	}
}

func WithValue(v string) func(*Parameter) {
	return func(p *Parameter) {
		p.Value = NewParameterValue(v, false)
	}
}

func WithSecret(v string) func(*Parameter) {
	return func(p *Parameter) {
		p.Value = NewParameterValue(v, true)
	}
}

func WithLastModified(lastModified time.Time) func(*Parameter) {
	return func(p *Parameter) {
		p.LastModified = lastModified
	}
}

func WithVersion(version int64) func(*Parameter) {
	return func(p *Parameter) {
		p.Version = version
	}
}

func AParameterName(builders ...func(*ParameterName)) *ParameterName {
	pn := &ParameterName{
		Capability:  "cap",
		Environment: "env",
		Application: "app",
		Key:         "key",
	}

	for _, builder := range builders {
		builder(pn)
	}

	return pn
}

func WithCapability(capability string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Capability = capability
	}
}

func WithEnvironment(environment string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Environment = environment
	}
}

func WithApplication(application string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Application = application
	}
}

func WithKey(key string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Key = key
	}
}
