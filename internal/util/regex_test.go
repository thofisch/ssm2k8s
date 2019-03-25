package util_test

import (
	"regexp"
	"testing"

	"github.com/thofisch/ssm2k8s/internal/assert"
	. "github.com/thofisch/ssm2k8s/internal/util"
)

func Test_findNamedGroups_with_named_group(t *testing.T) {
	groups := FindNamedGroups(regexp.MustCompile("(?P<group>.)"), "a")

	v, ok := groups["group"]

	assert.True(t, ok)
	assert.Equal(t, "a", v)

}

func Test_findNamedGroups_without_named_group(t *testing.T) {
	groups := FindNamedGroups(regexp.MustCompile("(.)"), "a")

	assert.Equal(t, 0, len(groups))
}

func Test_findNamedGroups_without_match(t *testing.T) {
	groups := FindNamedGroups(regexp.MustCompile("(?P<group>[a-z])"), "1")

	assert.Equal(t, 0, len(groups))
}
