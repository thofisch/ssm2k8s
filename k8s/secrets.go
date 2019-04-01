package k8s

import (
	"fmt"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
	"path/filepath"
	"runtime"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func PrintSecret() {
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "new-secret",
		},
		Data: map[string][]byte{"Key": []byte("Value")},
		Type: "Opaque",
	}

	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	s.Encode(secret, os.Stdout)

}

func GetSecrets() {
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
	secrets, err := clientset.CoreV1().Secrets("default").List(metav1.ListOptions{
		LabelSelector: "secrets.io/foo=bar",
	})
	if err != nil {
		panic(err)
	}
	for _, s := range secrets.Items {

		for k, v := range s.GetLabels() {
			fmt.Printf("l: %s = %s\n", k, v)
		}

		for k, v := range s.GetAnnotations() {
			fmt.Printf("a: %s = %s\n", k, v)
		}
		fmt.Println()

	}
	fmt.Println(time.Now().UTC().Format(time.RFC3339))
	//secret := &corev1.Secret{
	//	TypeMeta: metav1.TypeMeta{
	//		Kind:       "Secret",
	//		APIVersion: "v1",
	//	},
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name:        "new-secret",
	//		Labels:      map[string]string{"secrets.io/foo": "bar"},
	//		Annotations: map[string]string{"secrets.io/last-modified": time.Now().String()},
	//	},
	//	StringData: map[string]string{"key1": "foo", "key2": "bar"},
	//	//Data: map[string][]byte{"Key": []byte("Value")},
	//	Type: "Opaque",
	//}
	//
	//_, err = clientset.CoreV1().Secrets("default").Create(secret)
	//if err != nil {
	//	panic(err)
	//}
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
