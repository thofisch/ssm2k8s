// +build integration

package aws

import (
	"fmt"
	"testing"
)

func TestGetCallerIdentity(t *testing.T) {
	accountId, _ := GetAccountId()
	fmt.Println(accountId)
}
