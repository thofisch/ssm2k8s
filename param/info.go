package param

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ssm"
	"time"
)

type ParameterInfo struct {
	Name         *ParameterName
	Value        ParameterValue
	LastModified time.Time
	Version      int64
}

func (pi *ParameterInfo) String() string {
	return fmt.Sprintf("%s = %s", pi.Name, pi.Value)
}

func mapParameterInfo(p *ssm.Parameter) (*ParameterInfo, error) {
	name, err := parseParameterName(*p.Name)
	if err != nil {
		return nil, err
	}

	return &ParameterInfo{
		Name:         name,
		Value:        NewParameterValue(*p.Value, mapSecret(*p.Type)),
		LastModified: *p.LastModifiedDate,
		Version:      *p.Version,
	}, nil
}

func mapSecret(typeString string) bool {
	return ssm.ParameterTypeSecureString == typeString
}
