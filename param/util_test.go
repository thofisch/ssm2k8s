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

	assertEqual(t, 0, len(groups))
}

func Test_findNamedGroups_without_match(t *testing.T) {
	groups := findNamedGroups(regexp.MustCompile("(?P<group>[a-z])"), "1")

	assertEqual(t, 0, len(groups))
}
