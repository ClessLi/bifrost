package test

import (
	"testing"

	"github.com/yongPhone/bifrost/internal/pkg/password"
)

func TestPassword(t *testing.T) {
	pwd := password.Password("Bultgang")
	t.Log(pwd)
}
