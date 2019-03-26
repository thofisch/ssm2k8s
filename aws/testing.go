package aws

import "time"

func AParameter(builders ...func(*Parameter)) *Parameter {
	var pi = &Parameter{
		Name:         "",
		Value:        NewParameterValue("", false),
		Version:      int64(0),
		LastModified: time.Time{},
	}

	for _, builder := range builders {
		builder(pi)
	}

	return pi
}

func WithName(name string) func(*Parameter) {
	return func(p *Parameter) {
		p.Name = name
	}
}

func WithValue(v string) func(*Parameter) {
	return func(p *Parameter) {
		p.Value = NewParameterValue(v, false)
	}
}

func WithSecret(v string) func(*Parameter) {
	return func(p *Parameter) {
		p.Value = NewParameterValue(v, true)
	}
}

func WithLastModified(lastModified time.Time) func(*Parameter) {
	return func(p *Parameter) {
		p.LastModified = lastModified
	}
}

func WithVersion(version int64) func(*Parameter) {
	return func(p *Parameter) {
		p.Version = version
	}
}

//func AParameterName(builders ...func(*parameterName)) *parameterName {
//	pn := &parameterName{
//		Capability:  "cap",
//		Environment: "env",
//		Application: "app",
//		Key:         "key",
//	}
//
//	for _, builder := range builders {
//		builder(pn)
//	}
//
//	return pn
//}
//
//func WithCapability(capability string) func(*parameterName) {
//	return func(pn *parameterName) {
//		pn.Capability = capability
//	}
//}
//
//func WithEnvironment(environment string) func(*parameterName) {
//	return func(pn *parameterName) {
//		pn.Environment = environment
//	}
//}
//
//func WithApplication(application string) func(*parameterName) {
//	return func(pn *parameterName) {
//		pn.Application = application
//	}
//}
//
//func WithKey(key string) func(*parameterName) {
//	return func(pn *parameterName) {
//		pn.Key = key
//	}
//}
