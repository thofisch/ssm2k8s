package param

import (
	"testing"
)

func Test_parseParameterName_valid_name(t *testing.T) {
	pn, err := parseParameterName("/a/b/c/d")

	assertOk(t, err)
	assertEqual(t, "a", pn.Capability)
	assertEqual(t, "b", pn.Environment)
	assertEqual(t, "c", pn.Application)
	assertEqual(t, "d", pn.Key)
}

func Test_parseParameterName_invalid_name(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "only slashes", input: "////"},
		{name: "missing slash prefix", input: "a/b/c/d/e"},
		{name: "slash_cap", input: "/a"},
		{name: "slash_cap_slash", input: "/a/"},
		{name: "slash_cap_slash_env", input: "/a/b"},
		{name: "slash_cap_slash_env_slash", input: "/a/b/"},
		{name: "slash_cap_slash_env_slash_app", input: "/a/b/c"},
		{name: "slash_cap_slash_env_slash_app_slash", input: "/a/b/c/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash", input: "/a/b/c/d/"},
		{name: "slash_cap_slash_env_slash_app_slash_key_slash_extra", input: "/a/b/c/d/e"},
		{name: "slash_cap_slash_env_slash_slash_key", input: "/a/b//d"},
		{name: "slash_cap_slash_slash app_slash_key", input: "/a//c/d"},
		{name: "slash_slash_env_slash_app_slash_key", input: "//b/c/d"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseParameterName(test.input)

			assertNotOk(t, err)
		})
	}
}
