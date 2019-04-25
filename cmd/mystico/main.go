package main

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"strings"
	"time"
)

func main() {
	logger := logging.NewConsoleLogger()

	//var keyValuePairs keyValuePairs = make(keyValuePairs)
	//
	//flag.NewFlagSet()
	//
	//flag.Var(&keyValuePairs, "list1", "some description")
	//port := flag.Int("port", 8088, "Service Port Number")
	//flag.Parse()
	//
	//fmt.Printf("%#v\n", keyValuePairs)
	//fmt.Printf("%d\n", port)






}

//
//type keyValuePairs map[string]string
//
//func (kvp *keyValuePairs) String() string {
//	return fmt.Sprintf("Value: %#v\n", *kvp)
//}
//
//func (kvp *keyValuePairs) Set(value string) error {
//	split := strings.Split(value, "=")
//	if len(split) != 2 {
//		return fmt.Errorf("'%s' is unaccepted", value)
//	}
//
//	map[string]string(*kvp)[split[0]] = split[1]
//
//	return nil
//}
