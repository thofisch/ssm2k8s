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
	}
)

func NewSsmConfig(region string) *SsmConfig {
	return &SsmConfig{
		Region:    region,
		Recursive: true,
		Decrypt:   true,
	}
}

func NewSsmClient(config *SsmConfig) SsmClient {
	return &ssmClient{Config: config}
}

func (c *ssmClient) GetParametersByPath(path string) ([]*ssm.Parameter, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(c.Config.Region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	client := ssm.New(sess, aws.NewConfig().WithRegion(c.Config.Region))

	output, err := client.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(c.Config.Recursive),
		WithDecryption: aws.Bool(c.Config.Decrypt),
	})
	if err != nil {
		return nil, err
	}

	return output.Parameters, nil
}
