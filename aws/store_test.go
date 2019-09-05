package aws

import (
	"github.com/thofisch/ssm2k8s/internal/logging"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/internal/assert"
	"github.com/thofisch/ssm2k8s/internal/util"
)

func TestParameterStore_GetApplicationSecrets(t *testing.T) {
	first := time.Now().UTC()
	last := first.Add(time.Second)
	sut := NewParameterStoreWithClient(logging.NewNullLogger(), NewSsmClientStub(
		&ssm.Parameter{
			Name:             aws.String("/app/env/key1"),
			Value:            aws.String("secret"),
			Type:             aws.String(ssm.ParameterTypeSecureString),
			LastModifiedDate: aws.Time(first),
			Version:          aws.Int64(int64(0)),
		},
		&ssm.Parameter{
			Name:             aws.String("/app/env/key2"),
			Value:            aws.String("value"),
			Type:             aws.String(ssm.ParameterTypeString),
			LastModifiedDate: aws.Time(last),
			Version:          aws.Int64(int64(0)),
		},
	))

	parameters, err := sut.GetApplicationSecrets()

	assert.Ok(t, err)
	assert.Equal(t, domain.ApplicationSecrets{
		"env-app-secret": domain.ApplicationSecret{
			LastModified: last,
			Hash:         util.HashKeyValuePairs(map[string]string{"key1": "secret", "key2": "value"}),
			Data:         map[string]string{"key1": "secret", "key2": "value"},
		},
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

func (stub *SsmClientStub) PutParameter(name string, value string, overwrite bool) error {
	return nil
}

func (stub *SsmClientStub) DeleteParameter(name string) error {
	return nil
}

func Test_parseParameterName_invalid_name(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{name: "empty", key: ""},
		{name: "only slashes", key: "////"},
		{name: "missing slash prefix", key: "a/b/c/d/e"},
		{name: "slash_app", key: "/a"},
		{name: "slash_app_slash", key: "/a/"},
		{name: "slash_app_slash_env", key: "/a/b"},
		{name: "slash_app_slash_env_slash", key: "/a/b/"},
		{name: "slash_app_slash_env_slash_key_slash", key: "/a/b/c/"},
		{name: "slash_app_slash_env_slash_slash_key", key: "/a/b//d"},
		{name: "slash_app_slash_slash app_slash_key", key: "/a//c/d"},
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
	pn, err := parseParameterName("/a/b/c")

	assert.Ok(t, err)
	assert.Equal(t, "a", pn.Application)
	assert.Equal(t, "b", pn.Environment)
	assert.Equal(t, "c", pn.Key)
}

func Test_getSecretName(t *testing.T) {
	parameterName := parameterName{
		Application: "a",
		Environment: "b",
		Key:         "c",
	}

	assert.Equal(t, "b-a-secret", getSecretName(parameterName))
}
