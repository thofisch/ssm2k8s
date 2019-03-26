package ssm2k8s

import (
	"testing"

	"github.com/thofisch/ssm2k8s/internal/assert"
)

func Test_parseParameterName_invalid_name(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{name: "empty", key: ""},
		{name: "only slashes", key: "////"},
		{name: "missing slash prefix", key: "a/b/c/d/e"},
		{name: "slash_cap", key: "/a"},
		{name: "slash_cap_slash", key: "/a/"},
		{name: "slash_cap_slash_env", key: "/a/b"},
		{name: "slash_cap_slash_env_slash", key: "/a/b/"},
		{name: "slash_cap_slash_env_slash_app", key: "/a/b/c"},
		{name: "slash_cap_slash_env_slash_app_slash", key: "/a/b/c/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash", key: "/a/b/c/d/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash_extra", key: "/a/b/c/d/e"},
		{name: "slash_cap_slash_env_slash_slash_key", key: "/a/b//d"},
		{name: "slash_cap_slash_slash app_slash_key", key: "/a//c/d"},
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
	pn, err := parseParameterName("/a/b/c/d")

	assert.Ok(t, err)
	assert.Equal(t, "a", pn.Capability)
	assert.Equal(t, "b", pn.Environment)
	assert.Equal(t, "c", pn.Application)
	assert.Equal(t, "d", pn.Key)
}

//
//import (
//	"fmt"
//	. "github.com/thofisch/ssm2k8s/aws"
//	"testing"
//	"time"
//)
//
//func Test_aName(t *testing.T) {
//
//	var m = make(map[string]map[string]string)
//	if m != nil {
//
//	}
//
//	stub := NewParameterStoreStub(
//		aParameterInfo(
//			withName(aParameterName(
//				withApplication("foo"),
//				withKey("pghost"))),
//			withSecret("lala")),
//		aParameterInfo(
//			withName(aParameterName(
//				withApplication("foo"),
//				withKey("pguser"))),
//			withSecret("lala")),
//		aParameterInfo(
//			withName(aParameterName(
//				withApplication("foo"),
//				withKey("pgpassword"))),
//			withSecret("lala")),
//		aParameterInfo(
//			withName(aParameterName(
//				withApplication("foo"),
//				withKey("pgport"))),
//			withValue("1433")),
//		aParameterInfo(
//			withName(aParameterName(withKey("kafka-brokers"))),
//			withValue("broker1,broker2,broker3")),
//	)
//
//	parameters, _ := stub.GetParametersByPath("/p-project/")
//
//	result := NewResult(parameters)
//
//	for k, v := range result.Secrets {
//		fmt.Printf("Creating secret '\x1b[33m%s\x1b[0m' in namespace '\x1b[33m%s\x1b[0m' with the following values:\n", k, result.Name)
//		alignKeyValue(v)
//		fmt.Println()
//	}
//
//	for k, v := range result.Configs {
//		fmt.Printf("Creating config-map '\x1b[33m%s\x1b[0m' in namespace '\x1b[33m%s\x1b[0m' with the following values:\n", k, result.Name)
//		alignKeyValue(v)
//		fmt.Println()
//	}
//}
//
//func alignKeyValue(m map[string]ParameterValue) {
//	var max = 0
//	for k := range m {
//		l := len(k)
//		if l > max {
//			max = l
//		}
//	}
//
//	format := fmt.Sprintf("  \x1b[33m%%-%ds\x1b[0m = \x1b[34m%%s\x1b[0m\n", max)
//
//	for k, v := range m {
//		fmt.Printf(format, k, v)
//	}
//}
//
//type ParameterStoreStub struct {
//	parameterInfoList []*Parameter
//}
//
//func NewParameterStoreStub(parameters ...*Parameter) ssmClient {
//	return &ParameterStoreStub{parameterInfoList: parameters}
//}
//

//func (ps *ParameterStoreStub) GetParametersByPath(path string) ([]*Parameter, error) {
//	return ps.parameterInfoList, nil
//}
