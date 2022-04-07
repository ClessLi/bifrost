package password

import (

	// nolint:gosec
	"crypto/sha1"
	"fmt"
)

var Secret = "invisible_cloak" // 加盐
// var Secret = "" // 加盐

func Password(passwd string) string {
	sha1Inst := sha1.New() //nolint:gosec
	sha1Inst.Write([]byte(passwd))

	return fmt.Sprintf("%x", sha1Inst.Sum([]byte(Secret)))
}
