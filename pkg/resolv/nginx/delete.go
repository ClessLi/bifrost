package nginx

import (
	"os"
)

func Delete(config *Config) (err error) {
	l, err := config.List()
	if err != nil {
		return
	}

	for s := range l {
		err = os.Remove(s)
		if err != nil {
			return
		}
	}
	return nil
}
