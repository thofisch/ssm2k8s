package k8s

import (
	"io"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/thofisch/ssm2k8s/domain"
	"github.com/thofisch/ssm2k8s/internal/config"
	"github.com/thofisch/ssm2k8s/internal/logging"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
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
		Client Client
		Log    logging.Logger
	}
)

func NewSecretStore(log logging.Logger, namespace string) (SecretStore, error) {
	client, err := NewClient(log, Config{
		Namespace:     namespace,
		LabelSelector: LabelVersion,
	})
	if err != nil {
		return nil, err
	}

	return NewSecretStoreWithClient(log, client), nil
}

func NewSecretStoreWithClient(log logging.Logger, client Client) SecretStore {
	return &secretStore{
		Log:    log,
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
		version := s.GetLabels()[LabelVersion]

		err := verifyVersion(version)
		if err != nil {
			ss.Log.Debugf("SKIPPING: %q with version %q, due to %q \n", secretName, version, err)
			continue
		}

		Annotations := s.GetAnnotations()
		Data := getSecretData(s.Data)

		lastModified, _ := time.Parse(time.RFC3339, Annotations[AnnotationLastModified])
		secrets[secretName] = domain.ApplicationSecret{
			LastModified: lastModified,
			Hash:         Annotations[AnnotationHash],
			Data:         Data,
		}
	}

	return secrets, nil
}

func verifyVersion(version string) error {
	// TODO handle semantic version strategy
	_, err := semver.Make(strings.TrimPrefix("v", version))
	return err
}

func getSecretData(bytes map[string][]byte) domain.SecretData {
	result := make(domain.SecretData)

	for k, v := range bytes {
		result[k] = domain.DataSecret{
			Value:        string(v),
			Version:      0,
			LastModified: time.Time{},
		}
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
		StringData: getStringData(secret.Data),
		Type:       "Opaque",
	}
	return s
}

func getStringData(secretData domain.SecretData) map[string]string {
	result := make(map[string]string)

	for k, v := range secretData {
		result[k] = v.Value
	}

	return result
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
