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
	DefaultPollTimeout = 2 * time.Second
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
	namespace := "p-project"
	region := "eu-central-1"
	path := "/" + namespace

	fmt.Printf("#\n")
	fmt.Printf("# [SETUP] Pull from AWS SystemManager Parameters using: Path=\033[33m%s\033[0m, Account=\033[33m%s\033[0m, Region=\033[33m%s\033[0m\n", path, accountId, region)
	fmt.Printf("# [SETUP] Synchronize Kubernetes secrets in: Namespace=\033[33m%s\033[0m\n", namespace)
	fmt.Printf("#\n")

	config := ssm2k8s.Config{
		AccountId: accountId,
		Namespace: namespace,
		Region:    region,
	}

	secretStore, err := k8s.NewSecretStore(log, "default")
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
		sync:        ssm2k8s.NewSync(config, secretStore, parameterStore),
	}, nil
}

func (m *MainApp) run() {
	defer m.wg.Done()
	for {
		select {
		case <-time.After(m.PollTimeout):
			m.sync.SyncSecrets()

		case <-m.close:
			m.Log.Debug("Channel closed, quitting...")
			return
		}
	}
}
