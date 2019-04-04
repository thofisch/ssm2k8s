package util

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

func HashKeyValuePairs(kv map[string]string) string {
	keys := make([]string, 0, len(kv))

	for k := range kv {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	sb := strings.Builder{}

	for i, k := range keys {
		if i > 0 {
			fmt.Fprint(&sb, "&")
		}

		fmt.Fprintf(&sb, "%s=%s", k, kv[k])
	}

	hash := sha1.New()
	hash.Write([]byte(sb.String()))
	sum := hash.Sum(nil)

	return fmt.Sprintf("%x", sum)
}
