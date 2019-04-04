package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func GetAccountId() (string, error) {
	client := sts.New(session.Must(session.NewSession()))
	result, err := client.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return *result.Account, nil
}
