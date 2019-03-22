package param

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ParameterStore interface {
	GetParameters(path string) ([]*ParameterInfo, error)
}

type awsSystemManagerParameterStore struct {
	Region    string
	Recursive bool
	Decrypt   bool
}

func NewParameterStore(region string) ParameterStore {
	return &awsSystemManagerParameterStore{
		Region:    region,
		Recursive: true,
		Decrypt:   true,
	}
}

func (ps *awsSystemManagerParameterStore) GetParameters(path string) ([]*ParameterInfo, error) {
	parameters, err := ps.getParametersByPath(path)
	if err != nil {
		return nil, err
	}

	return mapParameters(parameters)
}

func (ps *awsSystemManagerParameterStore) getParametersByPath(path string) ([]*ssm.Parameter, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String(ps.Region)},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	client := ssm.New(sess, aws.NewConfig().WithRegion(ps.Region))

	output, err := client.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(ps.Recursive),
		WithDecryption: aws.Bool(ps.Decrypt),
	})
	if err != nil {
		return nil, err
	}

	return output.Parameters, nil
}

func mapParameters(parameters []*ssm.Parameter) ([]*ParameterInfo, error) {
	len := len(parameters)
	var parameterInfos = make([]*ParameterInfo, len)

	for i, p := range parameters {
		parameterInfo, err := mapParameterInfo(p)
		if err != nil {
			return nil, err
		}
		parameterInfos[i] = parameterInfo
	}

	return parameterInfos, nil
}
