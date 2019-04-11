package ssm2k8s

import (
	"github.com/thofisch/ssm2k8s/aws"
	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/internal/logging"
	"github.com/thofisch/ssm2k8s/k8s"
	"strings"
	"time"
)

type (
	Sync interface {
		SyncSecrets()
	}
	syncImpl struct {
		Log            logging.Logger
		Namespace      string
		SecretStore    k8s.SecretStore
		ParameterStore aws.ParameterStore
	}
)

func NewSync(logger logging.Logger, namespace string, secretStore k8s.SecretStore, parameterStore aws.ParameterStore) Sync {
	return &syncImpl{
		Log:            logger,
		SecretStore:    secretStore,
		ParameterStore: parameterStore,
		Namespace:      namespace,
	}
}

func (s *syncImpl) SyncSecrets() {
	// aws
	awsSecrets, err := s.ParameterStore.GetApplicationSecrets(s.Namespace)
	if err != nil {
		return
	}
	s.logApplicationSecrets(awsSecrets, "SSM AWS Parameters")

	// k8s
	k8sSecrets, err := s.SecretStore.GetApplicationSecrets()
	if err != nil {
		return
	}
	s.logApplicationSecrets(k8sSecrets, "Kubernetes")

	for secretName, secret := range awsSecrets {
		k8sSecret, ok := k8sSecrets[secretName]
		if !ok {
			continueOnError(s.SecretStore.CreateApplicationSecret(secret, secretName))
		} else {
			if secret.Hash == k8sSecret.Hash {
				s.Log.Debugf("Secret %q up-to-date (according to hash)", secretName)
				continue
			} else {
				continueOnError(s.SecretStore.UpdateApplicationSecret(secret, secretName))
			}
		}
	}

	for secretName := range k8sSecrets {
		_, ok := awsSecrets[secretName]

		if !ok {
			continueOnError(s.SecretStore.DeleteApplicationSecret(secretName))
		}
	}
}
func continueOnError(err error) bool {
	return err != nil
}

func (s *syncImpl) logApplicationSecrets(secrets domain.ApplicationSecrets, source string) {
	s.Log.Infof("Found %d secrets in %q", len(secrets), source)

	for appName, secret := range secrets {
		keys := make([]string, 0, len(secret.Data))

		for key := range secret.Data {
			keys = append(keys, key)
		}

		s.Log.Debugf("Name = %q, Hash = %q, LastModified = %q, Keys = %q",
			appName,
			secret.Hash,
			secret.LastModified.Format(time.RFC3339),
			strings.Join(keys, ", "),
		)
	}
}
