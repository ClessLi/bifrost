package test

import (
	"github.com/ClessLi/bifrost/internal/pkg/password"
	"testing"
)

func TestPassword(t *testing.T) {
	pwd := password.Password("Bultgang")
	t.Log(pwd)
}
