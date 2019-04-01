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
