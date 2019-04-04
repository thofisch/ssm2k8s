package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type (
	SsmClient interface {
		GetParametersByPath(path string) ([]*ssm.Parameter, error)
	}
	SsmConfig struct {
		Region    string
		Recursive bool
		Decrypt   bool
	}
	ssmClient struct {
		Config *SsmConfig
		ssm    *ssm.SSM
	}
)

func NewSsmConfig(region string) *SsmConfig {
	return &SsmConfig{
		Region:    region,
		Recursive: true,
		Decrypt:   true,
	}
}

func NewSsmClient(config *SsmConfig) (SsmClient, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region)},
	)
	if err != nil {
		return nil, err
	}

	ssm := ssm.New(session)

	return &ssmClient{
		Config: config,
		ssm:    ssm,
	}, nil
}

func (c *ssmClient) GetParametersByPath(path string) ([]*ssm.Parameter, error) {
	output, err := c.ssm.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(c.Config.Recursive),
		WithDecryption: aws.Bool(c.Config.Decrypt),
	})
	if err != nil {
		return nil, err
	}

	return output.Parameters, nil
}

func (c *ssmClient) PutParameter(name string, value string) error {
	_, err := c.ssm.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Overwrite: aws.Bool(true),
		Type:      aws.String(ssm.ParameterTypeSecureString),
	})

	return err
}
