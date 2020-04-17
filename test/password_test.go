package test

import (
	"github.com/ClessLi/go-nginx-conf-parser/internal/pkg/password"
	"testing"
)

func TestPassword(t *testing.T) {
	pwd := password.Password("ngadmin")
	t.Log(pwd)
}
