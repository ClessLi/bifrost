package resolv

import (
	"os"
	"os/exec"
)

func Check(config *Config, ng string) (err error) {
	err = config.Save()
	if err != nil {
		return
	}
	list, err := config.List()
	//fmt.Println("list", list)
	if err != nil {
		return
	}

	return check(ng, list[0])
}

func check(ng string, s string) error {
	cmd := exec.Command(ng, "-tc", s)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
