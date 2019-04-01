package main

import (
	"flag"
	"github.com/thofisch/ssm2k8s"
	"github.com/thofisch/ssm2k8s/k8s"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	stdlog "log"
)

func main() {

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
