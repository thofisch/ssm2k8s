package aws

import "time"

func AParameter(builders ...func(*Parameter)) *Parameter {
	var pi = &Parameter{
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

func WithName(pn *ParameterName) func(*Parameter) {
	return func(p *Parameter) {
		p.Name = pn
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

func AParameterName(builders ...func(*ParameterName)) *ParameterName {
	pn := &ParameterName{
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

func WithCapability(capability string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Capability = capability
	}
}

func WithEnvironment(environment string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Environment = environment
	}
}

func WithApplication(application string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Application = application
	}
}

func WithKey(key string) func(*ParameterName) {
	return func(pn *ParameterName) {
		pn.Key = key
	}
}
