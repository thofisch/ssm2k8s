package param

import (
	"regexp"
	"testing"
)

func Test_findNamedGroups_with_named_group(t *testing.T) {
	groups := findNamedGroups(regexp.MustCompile("(?P<group>.)"), "a")

	v, ok := groups["group"]

	assertTrue(t, ok)
	assertEqual(t, "a", v)

}

func Test_findNamedGroups_without_named_group(t *testing.T) {
	groups := findNamedGroups(regexp.MustCompile("(.)"), "a")

	assertEmpty(t, groups)
}

func Test_findNamedGroups_without_match(t *testing.T) {
	groups := findNamedGroups(regexp.MustCompile("(?P<group>[a-z])"), "1")

	assertEmpty(t, groups)
}
