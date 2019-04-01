package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

var wg sync.WaitGroup

var mainClose = make(chan bool)

func main() {
	//ssm2k8s.GenerateSecretManifests()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	go func() {
		for range signalChan {
			fmt.Println("\nReceived an interrupt, stopping services...")
			close(mainClose)
			os.Exit(0)
		}
	}()

	wg.Add(1)
	go runPoller()

	<-mainClose

	wg.Wait()
}

func runPoller() {
	defer wg.Done()
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Print(".")
		case _, ok := <-mainClose:
			fmt.Println(ok)
			return
		}
	}
}
