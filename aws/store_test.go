package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/thofisch/ssm2k8s/internal/assert"
	"testing"
	"time"
)

func TestParameterStore_GetParameters(t *testing.T) {
	then := time.Now().UTC()
	sut := NewParameterStoreWithClient(NewSsmClientStub(
		&ssm.Parameter{
			Name:             aws.String("/cap/env/app/key1"),
			Value:            aws.String("secret"),
			Type:             aws.String(ssm.ParameterTypeSecureString),
			LastModifiedDate: aws.Time(then),
			Version:          aws.Int64(int64(0)),
		},
		&ssm.Parameter{
			Name:             aws.String("/cap/env/app/key2"),
			Value:            aws.String("value"),
			Type:             aws.String(ssm.ParameterTypeString),
			LastModifiedDate: aws.Time(then),
			Version:          aws.Int64(int64(0)),
		},
	))

	parameters, err := sut.GetParameters("/path/")

	assert.Ok(t, err)
	assert.Equal(t, []*Parameter{
		AParameter(
			WithName("/cap/env/app/key1"),
			WithSecret("secret"),
			WithLastModified(then),
		),
		AParameter(
			WithName("/cap/env/app/key2"),
			WithValue("value"),
			WithLastModified(then),
		),
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
