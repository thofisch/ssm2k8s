package aws

type ParameterValue interface {
	GetValue() string
	IsSecret() bool
}

type parameterValue struct {
	value  string
	secret bool
}

func NewParameterValue(value string, secret bool) ParameterValue {
	return &parameterValue{value: value, secret: secret}
}

func (pv *parameterValue) GetValue() string {
	return pv.value
}

func (pv *parameterValue) IsSecret() bool {
	return pv.secret
}

func (pv *parameterValue) String() string {
	if pv.IsSecret() {
		return "*****"
	}
	return pv.GetValue()
}
