package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	kubeconfig := filepath.Join(
		homeDir(), ".kube", "nonprod",
	)

	fmt.Println(kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}

	secrets, err := clientset.CoreV1().Secrets("default").List(v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, s := range secrets.Items {

		fmt.Printf("%s\n", s.Name)

	}
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
