package main

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thofisch/ssm2k8s"
	"github.com/thofisch/ssm2k8s/k8s"
)

const (
	DefaultPollTimeout = 2 * time.Second
)

func main() {
	m, err := NewMain()
	if err != nil {
		panic(err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	go func() {
		for range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...")
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

	close chan bool
	wg    sync.WaitGroup

	sync ssm2k8s.Sync
}

func NewMain() (*MainApp, error) {

	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.Trace("TRACE")
	log.Debug("DEBUG")
	log.Info("INFO")
	log.Warn("WARN")
	//log.Error("ERROR")
	//log.Fatal("FATAL")


	log.Printf("#\n")
	log.Printf("# mysticod version 0.1 - synchronizing secrets\n")
	log.Printf("#\n")
	log.Printf("# Initializing... ")

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

	log.Printf("#\n")
	log.Printf("# [SETUP] Pull from AWS SystemManager Parameters using: Path=\033[33m%s\033[0m, Account=\033[33m%s\033[0m, Region=\033[33m%s\033[0m\n", path, accountId, region)
	log.Printf("# [SETUP] Synchronize Kubernetes secrets in: Namespace=\033[33m%s\033[0m\n", namespace)
	log.Printf("#\n")

	config := ssm2k8s.Config{
		AccountId: accountId,
		Namespace: namespace,
		Region:    region,
	}

	secretStore, err := k8s.NewSecretStore("default")
	if err != nil {
		return nil, err
	}

	parameterStore, err := aws.NewParameterStore(region)
	if err != nil {
		return nil, err
	}

	return &MainApp{
		PollTimeout: DefaultPollTimeout,
		close:       make(chan bool),
		sync:        ssm2k8s.NewSync(config, secretStore, parameterStore),
	}, nil
}

func (m *MainApp) run() {
	defer m.wg.Done()
	for {
		select {
		case <-time.After(m.PollTimeout):
			fmt.Print(".")

			m.sync.SyncSecrets()

		case _, ok := <-m.close:
			fmt.Println(ok)
			return
		}
	}
}
