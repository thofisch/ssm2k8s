package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/thofisch/ssm2k8s/internal/logging"
)

type (
	SsmClient interface {
		GetParametersByPath(path string) ([]*ssm.Parameter, error)
		PutParameter(name string, value string, overwrite bool) error
		DeleteParameter(name string) error
	}
	SsmConfig struct {
		Region    string
		Recursive bool
		Decrypt   bool
	}
	ssmClient struct {
		Log    logging.Logger
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

func NewSsmClient(logger logging.Logger, config *SsmConfig) (SsmClient, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region)},
	)
	if err != nil {
		return nil, err
	}

	ssm := ssm.New(session)

	return &ssmClient{
		Log:    logger,
		Config: config,
		ssm:    ssm,
	}, nil
}

func (c *ssmClient) GetParametersByPath(path string) ([]*ssm.Parameter, error) {
	var nextToken *string = nil
	var parameters []*ssm.Parameter

	for {
		c.Log.Debugf("aws ssm get-parameters-by-path(path=%q, recursive=%t, withDecryption=%t, nextToken=%q)",
			path, c.Config.Recursive, c.Config.Decrypt, safeString(nextToken))

		output, err := c.ssm.GetParametersByPath(&ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			Recursive:      aws.Bool(c.Config.Recursive),
			WithDecryption: aws.Bool(c.Config.Decrypt),
			NextToken:      nextToken,
		})
		if err != nil {
			c.Log.Errorf("ERROR: %s\n", err)
			return nil, err
		}

		c.Log.Debugf("Found %d parameters", len(output.Parameters))

		for _, p := range output.Parameters {
			parameters = append(parameters, p)
		}

		nextToken = output.NextToken

		if nextToken == nil {
			break
		}
		continue
	}

	return parameters, nil
}

func safeString(s *string) string {
	var ss string
	if s == nil {
		ss = ""
	} else {
		ss = *s
	}
	return ss
}

func (c *ssmClient) PutParameter(name string, value string, overwrite bool) error {
	_, err := c.ssm.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(name),
		Overwrite: aws.Bool(overwrite),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Value:     aws.String(value),
	})

	return err
}

func (c *ssmClient) DeleteParameter(name string) error {
	_, err := c.ssm.DeleteParameter(&ssm.DeleteParameterInput{
		Name: aws.String(name),
	})

	return err
}
