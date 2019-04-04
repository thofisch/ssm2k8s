package main

import (
	"flag"
	"fmt"
	"github.com/thofisch/ssm2k8s"
	"github.com/thofisch/ssm2k8s/k8s"
	"os"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	stdlog "log"
)

type keyValuePairs map[string]string

func (kvp *keyValuePairs) String() string {
	return fmt.Sprintf("Value: %#v\n", *kvp)
}

func (kvp *keyValuePairs) Set(value string) error {
	split := strings.Split(value, "=")
	if len(split) != 2 {
		return fmt.Errorf("'%s' is unaccepted", value)
	}

	map[string]string(*kvp)[split[0]] = split[1]

	return nil
}

func main() {
	var keyValuePairs keyValuePairs = make(keyValuePairs)

	flag.Var(&keyValuePairs, "list1", "some description")
	flag.Parse()

	fmt.Printf("%#v\n", keyValuePairs)

	logger := NewLogger("mysticod", "v0.1")

	level.Error(logger).Log("msg", "TCP Failure", "error", nil)

	port := flag.Int("port", 8088, "Service Port Number")
	flag.Parse()

	level.Info(logger).Log("msg", "Customer Service listening", "port", *port)

	ssm2k8s.CreateParameters()

	k8s.PrintSecret()
}

func NewLogger(serviceName string, version string) (logger log.Logger) {
	logger = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "svc", serviceName)
	logger = log.With(logger, "version", version)
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	return
}


//func CreateParameters() {
//	accountId, err := GetAccountId()
//	if err != nil {
//		panic(err)
//	}
//
//	path := "/p-project/"
//	fmt.Printf("Creating AWS SystemManager Parameters in \033[33m%s\033[0m from account \033[33m%s\033[0m \n", path, accountId)
//
//	newParameters := []struct {
//		environment string
//		application string
//		key         string
//		value       string
//		secret      bool
//	}{
//		{"prod", "default", "kafka_bootstrap_servers", "pkc-l9pve.eu-west-1.aws.confluent.cloud:9092", true},
//		{"prod", "default", "kafka_sasl_username", "2QIFZBFI32KT4I5Q", true},
//		{"prod", "default", "kafka_sasl_password", "1gOjGYMZrDqHQe3mlYpytZP5ScnQGWkJnwUimIX8KkCK+4GHfDOOl+DnBjOMrFre", true},
//	}
//
//	for _, v := range newParameters {
//		name := parameterName{"p-project", v.environment, v.application, v.key}
//
//		fmt.Printf("%s\n", &name)
//	}
//
//}
