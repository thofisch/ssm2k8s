package param_test

import (
	"fmt"
	. "github.com/thofisch/ssm2k8s/param"
	"testing"
	"time"
)

func Test_aName(t *testing.T) {

	stub := NewParameterStoreStub(
		aParameterInfo(
			withName(aParameterName(
				withApplication("foo"),
				withKey("pghost"))),
			withValue("lala"),
			isSecret()),
		aParameterInfo(
			withName(aParameterName(withKey("kafka-brokers"))),
			withValue("broker1,broker2,broker3")),
	)

	parameters, _ := stub.GetParameters("/p-project/")

	apps := make(map[string][]*ParameterInfo)

	for _, p := range parameters {
		list := apps[p.Name.Application]
		if list == nil {
			list = []*ParameterInfo{}
		}

		list = append(list, p)

		apps[p.Name.Application] = list

	}

	for k, v := range apps {

		fmt.Printf("Creating secret '%s' in namespace '%s' with the following values:\n", k, v[0].Name.Capability)

		for _, p := range v {
			fmt.Printf("  %s\n", p)
		}
		fmt.Println()
	}
}

type ParameterStoreStub struct {
	parameterInfoList []*ParameterInfo
}

func NewParameterStoreStub(parameters ...*ParameterInfo) ParameterStore {
	return &ParameterStoreStub{parameterInfoList: parameters}
}

func (ps *ParameterStoreStub) GetParameters(path string) ([]*ParameterInfo, error) {
	return ps.parameterInfoList, nil
}

func aParameterInfo(builders ...func(*ParameterInfo)) *ParameterInfo {
	pi := &ParameterInfo{
		Name:         aParameterName(),
		Secret:       No,
		Value:        "",
		Version:      0,
		LastModified: time.Time{},
	}

	for _, builder := range builders {
		builder(pi)
	}

	return pi
}

func withName(pn *ParameterName) func(*ParameterInfo) {
	return func(p *ParameterInfo) {
		p.Name = pn
	}
}

func withValue(v string) func(*ParameterInfo) {
	return func(p *ParameterInfo) {
		p.Value = v
	}
}

func isSecret() func(*ParameterInfo) {
	return func(pi *ParameterInfo) {
		pi.Secret = Yes
	}
}

func aParameterName(builders ...func(*ParameterName)) *ParameterName {
	pn := &ParameterName{
		Capability:  "selfservice",
		Environment: "prod",
		Application: "default",
		Key:         "",
	}

	for _, builder := range builders {
		builder(pn)
	}

	return pn
}

func withApplication(application string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Application = application
	}
}

func withKey(key string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Key = key
	}
}
