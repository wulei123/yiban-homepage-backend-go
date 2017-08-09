package yb

import (
	"testing"
	"../configuration"
)
func TestLogin(t *testing.T) {
	c := Login(configuration.Account,configuration.Password)
	if !CheckLogin(c) {
		t.Error("login failed")
	}
}