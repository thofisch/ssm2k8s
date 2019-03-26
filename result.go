package ssm2k8s

import . "github.com/thofisch/ssm2k8s/aws"

type ParameterValueMap map[string]ParameterValue

type N map[string]ParameterValueMap

type Result struct {
	Name    string
	Secrets N
	Configs N
}

//func NewResult(parameters []*Parameter) Result {
//	result := Result{Name: parameters[0].Name.Capability}
//	result.Secrets = make(N)
//	result.Configs = make(N)
//	secrets := fitlerSecrets(parameters, true)
//	configs := fitlerSecrets(parameters, false)
//	{
//		applications := applicationMap(secrets)
//		for a, pi := range applications {
//			keyValue := keyValueMap(pi)
//			result.Secrets[a] = keyValue
//		}
//	}
//	{
//		applications := applicationMap(configs)
//		for a, pi := range applications {
//			keyValue := keyValueMap(pi)
//			result.Configs[a] = keyValue
//		}
//	}
//	return result
//}
//
//func applicationMap(pi []*Parameter) map[string][]*Parameter {
//	m := make(map[string][]*Parameter)
//
//	for _, p := range pi {
//		m[p.Name.Application] = append(m[p.Name.Application], p)
//	}
//
//	return m
//}
//
//func keyValueMap(pi []*Parameter) ParameterValueMap {
//	m := make(ParameterValueMap)
//
//	for _, p := range pi {
//		m[p.Name.Key] = p.Value
//	}
//
//	return m
//}
//
//func fitlerSecrets(pi []*Parameter, b bool) []*Parameter {
//	var a []*Parameter
//
//	for _, p := range pi {
//		if p.Value.IsSecret() == b {
//			a = append(a, p)
//		}
//	}
//
//	return a
//}
