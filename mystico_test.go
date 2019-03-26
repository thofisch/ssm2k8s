package ssm2k8s

import (
	"testing"

	"github.com/thofisch/ssm2k8s/internal/assert"
)

func Test_parseParameterName_invalid_name(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{name: "empty", key: ""},
		{name: "only slashes", key: "////"},
		{name: "missing slash prefix", key: "a/b/c/d/e"},
		{name: "slash_cap", key: "/a"},
		{name: "slash_cap_slash", key: "/a/"},
		{name: "slash_cap_slash_env", key: "/a/b"},
		{name: "slash_cap_slash_env_slash", key: "/a/b/"},
		{name: "slash_cap_slash_env_slash_app", key: "/a/b/c"},
		{name: "slash_cap_slash_env_slash_app_slash", key: "/a/b/c/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash", key: "/a/b/c/d/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash_extra", key: "/a/b/c/d/e"},
		{name: "slash_cap_slash_env_slash_slash_key", key: "/a/b//d"},
		{name: "slash_cap_slash_slash app_slash_key", key: "/a//c/d"},
		{name: "slash_slash_env_slash_app_slash_key", key: "//b/c/d"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseParameterName(test.key)

			assert.NotOk(t, err)
		})
	}
}

func Test_parseParameterName_valid_name(t *testing.T) {
	pn, err := parseParameterName("/a/b/c/d")

	assert.Ok(t, err)
	assert.Equal(t, "a", pn.Capability)
	assert.Equal(t, "b", pn.Environment)
	assert.Equal(t, "c", pn.Application)
	assert.Equal(t, "d", pn.Key)
}
