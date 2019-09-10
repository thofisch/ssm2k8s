package domain

import (
	"time"
)

type (
	ApplicationSecret struct {
		Path         string
		LastModified time.Time
		Hash         string
		Data         SecretData
	}

	ApplicationSecrets map[string]ApplicationSecret
	SecretData         map[string]DataSecret

	DataSecret struct {
		Value        string
		Version      int64
		LastModified time.Time
	}
)

//func (secrets ApplicationSecrets) String() string {
//	sb := &strings.Builder{}
//
//	_, err := fmt.Fprintf(sb, "Application Secrets:\n")
//	if err != nil {
//		return ""
//	}
//
//	for appName, secret := range secrets {
//		keys := make([]string, 0, len(secret.Data))
//
//		for key := range secret.Data {
//			keys = append(keys, fmt.Sprintf("\033[33m%s\033[0m", key))
//		}
//
//		_, err = fmt.Fprintf(sb, " name = \033[33m%s\033[0m, hash = \033[33m%s\033[0m, last-modified = \033[33m%s\033[0m, keys = %s\n",
//			appName,
//			secret.Hash,
//			secret.LastModified.Format(time.RFC3339),
//			strings.Join(keys, ", "),
//		)
//		if err != nil {
//			return ""
//		}
//	}
//
//	return sb.String()
//}
