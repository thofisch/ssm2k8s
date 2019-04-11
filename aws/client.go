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

	//log.NewContextLogger(logrus.Fields{
	//	"path":      path,
	//	"recursive": c.Config.Recursive,
	//	"decrypt":   c.Config.Decrypt,
	//	"region":    c.Config.Region,
	//})

	for {
		c.Log.Debugf("aws ssm get-parameters-by-path(path=%q, recursive=%t, withDecryption=%t, nextToken=%s\n",
			path, c.Config.Recursive, c.Config.Decrypt, nextToken)

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

		c.Log.Debugf("Found %d parameters\n", len(output.Parameters))

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

func (c *ssmClient) PutParameter(name string, value string) error {
	_, err := c.ssm.PutParameter(&ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Overwrite: aws.Bool(true),
		Type:      aws.String(ssm.ParameterTypeSecureString),
	})

	return err
}
