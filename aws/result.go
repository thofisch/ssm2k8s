package aws

type ParameterValueMap map[string]ParameterValue

type N map[string]ParameterValueMap

type Result struct {
	Name    string
	Secrets N
	Configs N
}

func NewResult(parameters []*parameter) Result {
	result := Result{Name: parameters[0].Name.Capability}
	result.Secrets = make(N)
	result.Configs = make(N)
	secrets := fitlerSecrets(parameters, true)
	configs := fitlerSecrets(parameters, false)
	{
		applications := applicationMap(secrets)
		for a, pi := range applications {
			keyValue := keyValueMap(pi)
			result.Secrets[a] = keyValue
		}
	}
	{
		applications := applicationMap(configs)
		for a, pi := range applications {
			keyValue := keyValueMap(pi)
			result.Configs[a] = keyValue
		}
	}
	return result
}

func applicationMap(pi []*parameter) map[string][]*parameter {
	m := make(map[string][]*parameter)

	for _, p := range pi {
		m[p.Name.Application] = append(m[p.Name.Application], p)
	}

	return m
}

func keyValueMap(pi []*parameter) ParameterValueMap {
	m := make(ParameterValueMap)

	for _, p := range pi {
		m[p.Name.Key] = p.Value
	}

	return m
}

func fitlerSecrets(pi []*parameter, b bool) []*parameter {
	var a []*parameter

	for _, p := range pi {
		if p.Value.IsSecret() == b {
			a = append(a, p)
		}
	}

	return a
}
