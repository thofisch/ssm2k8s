package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	"runtime"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	//sess, err := session.NewSessionWithOptions(session.Options{
	//	Config:            aws.Config{Region: aws.String("eu-central-1")},
	//	SharedConfigState: session.SharedConfigEnable,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//ssmsvc := ssm.New(sess, aws.NewConfig().WithRegion("eu-central-1"))
	//keyname := "/p-project/"
	//withDecryption := true
	//recursive := true
	//
	//// ssm.GetParametersByPath()
	//
	//output, err := ssmsvc.GetParametersByPath(&ssm.GetParametersByPathInput{
	//	Path:           &keyname,
	//	Recursive:      &recursive,
	//	WithDecryption: &withDecryption,
	//})
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, p := range output.Parameters {
	//	fmt.Printf("%s = %s\n", *p.Name, *p.Value)
	//}

	kubeconfig := filepath.Join(
		homeDir(), ".kube", "np",
	)

	fmt.Println(kubeconfig)

	//clientcmd.

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
