package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/thofisch/ssm2k8s"
	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"github.com/thofisch/ssm2k8s/k8s"
)

const (
	DefaultPollTimeout         = 30
	KubernetesNamespaceEnvName = "KUBERNETES_NAMESPACE"
)

func main() {
	logger := logging.NewLogger()

	m, err := NewMain(logger)
	if err != nil {
		panic(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	go func() {
		for range signalChan {
			logger.Debug("Received shutdown signal, terminating...")
			close(m.close)
		}
	}()

	m.wg.Add(1)
	go m.run()

	<-m.close
	m.wg.Wait()
}

type MainApp struct {
	PollTimeout time.Duration

	Log logging.Logger

	close chan bool
	wg    sync.WaitGroup

	sync ssm2k8s.Sync
}

func NewMain(log logging.Logger) (*MainApp, error) {
	fmt.Printf("#\n")
	fmt.Printf("# mysticod version 0.1 - synchronizing secrets\n")
	fmt.Printf("#\n")
	fmt.Printf("# Initializing... ")

	accountId, err := aws.GetAccountId()
	if err != nil {
		fmt.Printf("\n# [ERROR] %s", err)
		return nil, err
	} else {
		fmt.Printf("OK\n")
	}

	// TODO -- get namespace from configuration

	namespace, ok := os.LookupEnv(KubernetesNamespaceEnvName)
	if !ok {
		panic(fmt.Sprintf("Expected environment variable %q to containing namespace information", KubernetesNamespaceEnvName))
	}

	region := "eu-central-1"

	fmt.Printf("#\n")
	printConfig(map[string]string{
		"config.aws.region":              region,
		"config.aws.accountId":           accountId,
		"config.aws.ssm.path":            "/" + namespace,
		"config.aws.ssm.recursive":       "true",
		"config.aws.ssm.with_decryption": "true",
		"config.kubernetes.namespace":    namespace,
	})
	fmt.Printf("#\n")

	secretStore, err := k8s.NewSecretStore(log, namespace)
	if err != nil {
		return nil, err
	}

	parameterStore, err := aws.NewParameterStore(log, region)
	if err != nil {
		return nil, err
	}

	return &MainApp{
		PollTimeout: DefaultPollTimeout,
		Log:         log,
		close:       make(chan bool),
		sync:        ssm2k8s.NewSync(log, namespace, secretStore, parameterStore),
	}, nil
}

func printConfig(config map[string]string) {
	max := 0
	for k := range config {
		l := len(k)
		if l > max {
			max = l
		}
	}

	format := fmt.Sprintf("# \033[34m%%-%ds\033[0m \033[33m%%s\033[0m\n", max)

	for k, v := range config {
		fmt.Printf(format, k, v)
	}

}

func (m *MainApp) run() {
	defer m.wg.Done()
	for {
		select {
		case <-time.After(m.PollTimeout * time.Second):
			m.Log.Info("Synchronizing secrets")
			m.sync.SyncSecrets()
			m.Log.Infof("Done synchronizing secrets. Waiting %d seconds", m.PollTimeout)

		case <-m.close:
			m.Log.Debug("Channel closed, quitting...")
			return
		}
	}
}
