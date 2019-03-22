package param

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ssm"
	"time"
)

type ParameterInfo struct {
	Name         *ParameterName
	Value        *ParameterValue
	LastModified time.Time
	Version      int64
}

func (pi *ParameterInfo) String() string {
	return fmt.Sprintf("%s = %s", pi.Name, *pi.Value)
}

func mapParameterInfo(p *ssm.Parameter) (*ParameterInfo, error) {
	name, err := parseParameterName(*p.Name)
	if err != nil {
		return nil, err
	}

	var value = NewParameterValue(*p.Value, mapSecret(*p.Type))

	pi := &ParameterInfo{}
	pi.Name = name
	pi.Value = &value
	pi.LastModified = *p.LastModifiedDate
	pi.Version = *p.Version

	return pi, nil
}

func mapSecret(typeString string) bool {
	return ssm.ParameterTypeSecureString == typeString
}
