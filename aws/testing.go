package aws

import "time"

func AParameter(builders ...func(*parameter)) *parameter {
	var pi = &parameter{
		Name:         AParameterName(),
		Value:        NewParameterValue("", false),
		Version:      0,
		LastModified: time.Time{},
	}

	for _, builder := range builders {
		builder(pi)
	}

	return pi
}

func WithName(pn *parameterName) func(*parameter) {
	return func(p *parameter) {
		p.Name = pn
	}
}

func WithValue(v string) func(*parameter) {
	return func(p *parameter) {
		p.Value = NewParameterValue(v, false)
	}
}

func WithSecret(v string) func(*parameter) {
	return func(p *parameter) {
		p.Value = NewParameterValue(v, true)
	}
}

func WithLastModified(lastModified time.Time) func(*parameter) {
	return func(p *parameter) {
		p.LastModified = lastModified
	}
}

func WithVersion(version int64) func(*parameter) {
	return func(p *parameter) {
		p.Version = version
	}
}

func AParameterName(builders ...func(*parameterName)) *parameterName {
	pn := &parameterName{
		Capability:  "cap",
		Environment: "env",
		Application: "app",
		Key:         "key",
	}

	for _, builder := range builders {
		builder(pn)
	}

	return pn
}

func WithCapability(capability string) func(*parameterName) {
	return func(pn *parameterName) {
		pn.Capability = capability
	}
}

func WithEnvironment(environment string) func(*parameterName) {
	return func(pn *parameterName) {
		pn.Environment = environment
	}
}

func WithApplication(application string) func(*parameterName) {
	return func(pn *parameterName) {
		pn.Application = application
	}
}

func WithKey(key string) func(*parameterName) {
	return func(pn *parameterName) {
		pn.Key = key
	}
}
