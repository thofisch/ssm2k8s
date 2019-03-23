package param

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"testing"
	"time"
)

func Test_mapParameterInfo_can_map_values(t *testing.T) {
	value := "val"
	typeString := ssm.ParameterTypeSecureString
	lastModified, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00")
	version := int64(1)
	parameter := &ssm.Parameter{
		Name:             aws.String("/a/b/c/d"),
		Value:            aws.String(value),
		Version:          aws.Int64(version),
		Type:             aws.String(typeString),
		LastModifiedDate: aws.Time(lastModified),
	}

	result, err := mapParameterInfo(parameter)

	assertOk(t, err)
	assertEqual(t, "a", result.Name.Capability)
	assertEqual(t, "b", result.Name.Environment)
	assertEqual(t, "c", result.Name.Application)
	assertEqual(t, "d", result.Name.Key)
	assertEqual(t, value, result.Value.GetValue())
	//assertEqual(t, Yes, result.Secret)
	assertEqual(t, lastModified, result.LastModified)
	assertEqual(t, version, result.Version)
}
