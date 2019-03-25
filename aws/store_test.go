package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/thofisch/ssm2k8s/internal/assert"
)

func TestParameterStore_GetParameters(t *testing.T) {
	then := time.Now().UTC()
	sut := NewParameterStoreWithClient(NewSsmClientStub(
		&ssm.Parameter{
			Version:          aws.Int64(1),
			LastModifiedDate: aws.Time(then),
			Type:             aws.String(ssm.ParameterTypeSecureString),
			Value:            aws.String("val"),
			Name:             aws.String("/cap/env/app/key"),
		},
		&ssm.Parameter{
			Version:          aws.Int64(1),
			LastModifiedDate: aws.Time(then),
			Type:             aws.String(ssm.ParameterTypeString),
			Value:            aws.String("val"),
			Name:             aws.String("/cap/env/app/key"),
		},
	))

	parameters, err := sut.GetParameters("/path/")

	assert.Ok(t, err)
	assert.Equal(t, []*parameter{
		AParameter(
			WithName(AParameterName()),
			WithVersion(1),
			WithSecret("val"),
			WithLastModified(then)),
		AParameter(
			WithName(AParameterName()),
			WithVersion(1),
			WithValue("val"),
			WithLastModified(then)),
	}, parameters)
}

type SsmClientStub struct {
	Parameters []*ssm.Parameter
}

func NewSsmClientStub(parameters ...*ssm.Parameter) *SsmClientStub {
	return &SsmClientStub{Parameters: parameters}
}

func (stub *SsmClientStub) GetParametersByPath(path string) ([]*ssm.Parameter, error) {
	return stub.Parameters, nil
}

func Test_parseParameterName_invalid_name(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{name: "empty", key: ""},
		{name: "only slashes", key: "////"},
		{name: "missing slash prefix", key: "a/b/c/d/e"},
		{name: "slash_cap", key: "/a"},
		{name: "slash_cap_slash", key: "/a/"},
		{name: "slash_cap_slash_env", key: "/a/b"},
		{name: "slash_cap_slash_env_slash", key: "/a/b/"},
		{name: "slash_cap_slash_env_slash_app", key: "/a/b/c"},
		{name: "slash_cap_slash_env_slash_app_slash", key: "/a/b/c/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash", key: "/a/b/c/d/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash_extra", key: "/a/b/c/d/e"},
		{name: "slash_cap_slash_env_slash_slash_key", key: "/a/b//d"},
		{name: "slash_cap_slash_slash app_slash_key", key: "/a//c/d"},
		{name: "slash_slash_env_slash_app_slash_key", key: "//b/c/d"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseParameterName(test.key)

			assert.NotOk(t, err)
		})
	}
}

func Test_parseParameterName_valid_name(t *testing.T) {
	pn, err := parseParameterName("/a/b/c/d")

	assert.Ok(t, err)
	assert.Equal(t, "a", pn.Capability)
	assert.Equal(t, "b", pn.Environment)
	assert.Equal(t, "c", pn.Application)
	assert.Equal(t, "d", pn.Key)
}
