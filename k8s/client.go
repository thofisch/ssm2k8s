package k8s

import (
	"os"

	"github.com/thofisch/ssm2k8s/internal/logging"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client interface {
	GetSecrets() ([]coreV1.Secret, error)
	CreateSecret(*coreV1.Secret) error
	UpdateSecret(*coreV1.Secret) error
	DeleteSecret(secretName string) error
}

type client struct {
	Log       logging.Logger
	Config    Config
	Clientset *kubernetes.Clientset
}

type Config struct {
	KubeconfigPath string
	CurrentContext string
	Namespace      string
	LabelSelector  string
}

func NewClient(logger logging.Logger, c Config) (Client, error) {

	config, err := getConfig(c)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &client{
		Log:       logger,
		Config:    c,
		Clientset: clientset,
	}, nil
}

func getConfig(c Config) (*rest.Config, error) {
	// are we inCluster?
	_, exist := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if exist {
		return rest.InClusterConfig()
	}

	// did we supply a kubeconfigPath?
	if len(c.KubeconfigPath) > 0 {
		return clientcmd.BuildConfigFromFlags("", c.KubeconfigPath)
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: c.CurrentContext}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides).ClientConfig()
}

func (c *client) GetSecrets() ([]coreV1.Secret, error) {
	c.Log.Debugf("Listing secrets in namespace %q", c.Config.Namespace)
	secretList, err := c.Clientset.CoreV1().Secrets(c.Config.Namespace).List(metaV1.ListOptions{
		LabelSelector: c.Config.LabelSelector,
	})
	if err != nil {
		c.Log.Errorf("ERROR: %s", err)
		return nil, err
	}

	secrets := make([]coreV1.Secret, len(secretList.Items))

	for i, s := range secretList.Items {
		secrets[i] = s
	}

	return secrets, nil
}

func (c *client) CreateSecret(secret *coreV1.Secret) error {
	c.Log.Debugf("Creating secret %q in namespace %q", secret.Name, c.Config.Namespace)
	_, err := c.Clientset.CoreV1().Secrets(c.Config.Namespace).Create(secret)
	if err != nil {
		c.Log.Errorf("ERROR: %s", err)
	}

	return err
}

func (c *client) UpdateSecret(secret *coreV1.Secret) error {
	c.Log.Debugf("Updating secret %q in namespace %q", secret.Name, c.Config.Namespace)
	_, err := c.Clientset.CoreV1().Secrets(c.Config.Namespace).Update(secret)
	if err != nil {
		c.Log.Errorf("ERROR: %s", err)
	}

	return err
}

func (c *client) DeleteSecret(name string) error {
	c.Log.Debugf("Deleting secret %q in namespace %n", name, c.Config.Namespace)
	err := c.Clientset.CoreV1().Secrets(c.Config.Namespace).Delete(name, &metaV1.DeleteOptions{})
	if err != nil {
		c.Log.Errorf("ERROR: %s", err)
	}

	return err
}
