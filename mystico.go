package ssm2k8s

import (
	"fmt"
	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/k8s"
)

type Sync interface {
	SyncSecrets()
}

type syncImpl struct {
	secretStore    k8s.SecretStore
	parameterStore aws.ParameterStore

	accountId string
	region    string
	path      string
	namespace string
}

type Config struct {
	AccountId string
	Region    string
	Namespace string
}

func NewSync(config Config, secretStore k8s.SecretStore, parameterStore aws.ParameterStore) Sync {
	return &syncImpl{
		secretStore:    secretStore,
		parameterStore: parameterStore,
		accountId:      config.AccountId,
		namespace:      "default", //config.Namespace,
		region:         config.Region,
		path:           "/" + config.Namespace,
	}
}

func (s *syncImpl) SyncSecrets() {
	// aws
	fmt.Printf("pulling aws ssm parameters\n")
	secrets, err := s.parameterStore.GetApplicationSecrets(s.path)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		return
	}

	domain.PrintApplicationSecrets(secrets, "AWS SSM: "+s.path)

	// k8s
	fmt.Printf("synchornizing k8s secrets\n")

	k8sSecrets, err := s.secretStore.GetApplicationSecrets()
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
		return
	}

	domain.PrintApplicationSecrets(k8sSecrets, "Namespace: "+s.namespace)

	for secretName, secret := range secrets {
		k8sSecret, ok := k8sSecrets[secretName]
		if !ok {

			fmt.Printf("Creating new secret \033[33m%s\033[0m in \033[33m%s\033[0m... ", secretName, s.namespace)

			err := s.secretStore.CreateApplicationSecret(secret, secretName)
			if err == nil {
				fmt.Printf("OK\n")
			} else {
				fmt.Printf("ERROR!\n%s\n", err)
			}
		} else {
			fmt.Printf("Found secret \033[33m%s\033[0m in \033[33m%s\033[0m, checking hash... ", secretName, s.namespace)

			if secret.Hash == k8sSecret.Hash {
				fmt.Printf("\033[32mOK\033[0m\n")
				continue
			} else {
				fmt.Printf("\033[34mWARN\033[0m\n")
			}

			err := s.secretStore.UpdateApplicationSecret(secret, secretName)
			if err == nil {
				fmt.Printf("OK\n")
			} else {
				fmt.Printf("ERROR!\n%s\n", err)
			}
		}
	}

	for secretName := range k8sSecrets {
		_, ok := secrets[secretName]

		if !ok {
			fmt.Printf("Deleting secret \033[33m%s\033[0m in \033[33m%s\033[0m, checking hash... ", secretName, s.namespace)

			s.secretStore.DeleteApplicationSecret(secretName)
		}
	}
}
