package resolv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func Backup(config *Config, name string) (path string, err error) {
	if name != "" {
		if !isTgz(name) {
			name = strings.TrimSpace(name)
			name = strings.TrimRight(name, ".")
		}
	} else {
		dt := time.Now().Format("2006-01-02_150405")
		list, err := config.List()
		if err != nil {
			return "", err
		}
		dir := filepath.Dir(list[0])
		name = filepath.Join(dir, fmt.Sprintf("nginx.conf.bak%s.tgz", dt))
	}
	list, err := config.List()
	if err != nil {
		return "", err
	}
	path = filepath.Join(list[0], name)
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if _, stat := os.Stat(path); os.IsExist(stat) {
		err = stat
		return
	}

	err = tgz(path, list)
	return
}

func tgz(path string, list []string) (err error) {
	tar, err := exec.LookPath("tar")
	if err != nil {
		return
	}
	args := append([]string{"-zcf", path}, list...)
	cmd := exec.Command(tar, args...)
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}

func isTgz(name string) bool {
	reg := regexp.MustCompile(`^.*\.tar\.gz|\.tgz$`)
	return reg.MatchString(name)
}
