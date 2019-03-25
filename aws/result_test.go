package aws

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
//	parameterInfoList []*parameter
//}
//
//func NewParameterStoreStub(parameters ...*parameter) ssmClient {
//	return &ParameterStoreStub{parameterInfoList: parameters}
//}
//

//func (ps *ParameterStoreStub) GetParametersByPath(path string) ([]*parameter, error) {
//	return ps.parameterInfoList, nil
//}
