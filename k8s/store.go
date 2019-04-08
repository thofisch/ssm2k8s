package k8s

import (
	"github.com/thofisch/ssm2k8s/internal/logging"
	"io"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
	"time"

	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/internal/config"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ApplicationPrefix      = "mystico.io"
	LabelVersion           = ApplicationPrefix + "/version"
	AnnotationLastModified = ApplicationPrefix + "/last-modified"
	AnnotationHash         = ApplicationPrefix + "/hash"
)

type (
	SecretStore interface {
		GetApplicationSecrets() (domain.ApplicationSecrets, error)
		CreateApplicationSecret(secret domain.ApplicationSecret, secretName string) error
		UpdateApplicationSecret(secret domain.ApplicationSecret, secretName string) error
		DeleteApplicationSecret(secretName string) error
	}
	secretStore struct {
		Log    logging.Logger
		Client Client
	}
)

func NewSecretStore(logger logging.Logger, namespace string) (SecretStore, error) {
	client, err := NewClient(logger, Config{
		Namespace:     namespace,
		LabelSelector: LabelVersion,
	})
	if err != nil {
		return nil, err
	}

	return NewSecretStoreWithClient(logger, client), nil
}

func NewSecretStoreWithClient(logger logging.Logger, client Client) SecretStore {
	return &secretStore{
		Log:    logger,
		Client: client,
	}
}

func (ss *secretStore) GetApplicationSecrets() (domain.ApplicationSecrets, error) {
	k8sSecrets, err := ss.Client.GetSecrets()
	if err != nil {
		return nil, err
	}

	secrets := make(domain.ApplicationSecrets)

	for _, s := range k8sSecrets {
		secretName := s.GetName()
		Annotations := s.GetAnnotations()
		Data := getKeyValueMap(s.Data)

		lastModified, _ := time.Parse(time.RFC3339, Annotations[AnnotationLastModified])
		secrets[secretName] = domain.ApplicationSecret{
			LastModified: lastModified,
			Hash:         Annotations[AnnotationHash],
			Data:         Data,
		}
	}

	return secrets, nil
}

func getKeyValueMap(bytes map[string][]byte) map[string]string {
	result := make(map[string]string)

	for k, v := range bytes {
		result[k] = string(v)
	}

	return result
}

func (ss *secretStore) CreateApplicationSecret(secret domain.ApplicationSecret, secretName string) error {
	s := ss.createSecret(secret, secretName)

	//return nil

	//return printSecret(s, os.Stderr)

	return ss.Client.CreateSecret(s)
}

func (ss *secretStore) createSecret(secret domain.ApplicationSecret, secretName string) *coreV1.Secret {
	s := &coreV1.Secret{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: secretName,
			Labels: map[string]string{
				LabelVersion: config.Version,
			},
			Annotations: map[string]string{
				AnnotationLastModified: secret.LastModified.Format(time.RFC3339),
				AnnotationHash:         secret.Hash},
		},
		StringData: secret.Data,
		Type:       "Opaque",
	}
	return s
}

func printSecret(secret *coreV1.Secret, w io.Writer) error {
	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

	s.Encode(secret, w)

	return nil
}

func (ss *secretStore) UpdateApplicationSecret(secret domain.ApplicationSecret, secretName string) error {
	s := ss.createSecret(secret, secretName)

	//return nil

	//return printSecret(s, os.Stderr)

	return ss.Client.UpdateSecret(s)
}

func (ss *secretStore) DeleteApplicationSecret(secretName string) error {
	return ss.Client.DeleteSecret(secretName)
}
