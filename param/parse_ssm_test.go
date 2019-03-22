package param

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"testing"
	"time"
)

func Test_parser(t *testing.T) {
	parameters, err := parseParameterName("/a/b/c/d")

	assertOk(t, err)
	assertEqual(t, "a", parameters.Capability)
	assertEqual(t, "b", parameters.Environment)
	assertEqual(t, "c", parameters.Application)
	assertEqual(t, "d", parameters.Key)
}

func Test_parseParameterName(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "only slashes", input: "////"},
		{name: "missing slash prefix", input: "a/b/c/d/e"},
		{name: "slash_cap", input: "/a"},
		{name: "slash_cap_slash", input: "/a/"},
		{name: "slash_cap_slash_env", input: "/a/b"},
		{name: "slash_cap_slash_env_slash", input: "/a/b/"},
		{name: "slash_cap_slash_env_slash_app", input: "/a/b/c"},
		{name: "slash_cap_slash_env_slash_app_slash", input: "/a/b/c/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash", input: "/a/b/c/d/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash_extra", input: "/a/b/c/d/e"},
		{name: "slash_cap_slash_env_slash_slash_key", input: "/a/b//d"},
		{name: "slash_cap_slash_slash app_slash_key", input: "/a//c/d"},
		{name: "slash_slash_env_slash_app_slash_key", input: "//b/c/d"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseParameterName(test.input)

			assertNotOk(t, err)
		})
	}
}

func Test_toParameterInfo_can_map_values(t *testing.T) {
	name := "/a/b/c/d"
	value := "val"
	typeString := "SecureString"
	lastModified, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00")
	version := int64(1)

	parameter := aParameter(
		withName(name),
		withValue(value),
		withType(typeString),
		withLastModifiedDate(lastModified),
		withVersion(version),
	)

	result, err := toParameterInfo(parameter)

	assertOk(t, err)
	assertEqual(t, "a", result.Name.Capability)
	assertEqual(t, "b", result.Name.Environment)
	assertEqual(t, "c", result.Name.Application)
	assertEqual(t, "d", result.Name.Key)
	assertEqual(t, value, result.Value)
	assertEqual(t, typeString, result.Type)
	assertEqual(t, lastModified, result.LastModified)
	assertEqual(t, version, result.Version)
}

func (pn ParameterName) String() string {
	return fmt.Sprintf("/%s/%s/%s/%s", pn.Capability, pn.Environment, pn.Application, pn.Key)
}

func TestName(t *testing.T) {

	parameters := NewParameterStore("eu-central-1").GetParameters("/p-project/")

	for _, p := range parameters {
		fmt.Printf("%-50s = '%s' [%s]\n", p.Name, p.Value, p.Type)

	}
}

/**********************************************************************************************************************
 ***
 *** Test Data Builder
 ***
 **********************************************************************************************************************/

func aParameter(builders ...func(*ssm.Parameter)) *ssm.Parameter {
	node := &ssm.Parameter{
		Name:             aws.String(""),
		Value:            aws.String(""),
		Version:          aws.Int64(0),
		Type:             aws.String(""),
		LastModifiedDate: aws.Time(time.Time{}),
	}

	for _, build := range builders {
		build(node)
	}

	return node
}

func withName(n string) func(parameter *ssm.Parameter) {
	return func(p *ssm.Parameter) {
		p.Name = aws.String(n)
	}
}

func withValue(v string) func(parameter *ssm.Parameter) {
	return func(p *ssm.Parameter) {
		p.Value = aws.String(v)
	}
}

func withType(t string) func(parameter *ssm.Parameter) {
	return func(p *ssm.Parameter) {
		p.Type = aws.String(t)
	}
}

func withLastModifiedDate(t time.Time) func(parameter *ssm.Parameter) {
	return func(p *ssm.Parameter) {
		p.LastModifiedDate = aws.Time(t)
	}
}
func withVersion(v int64) func(parameter *ssm.Parameter) {
	return func(p *ssm.Parameter) {
		p.Version = aws.Int64(v)
	}
}
