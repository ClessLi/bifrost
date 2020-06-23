package nginx

import (
	"os"
)

func Delete(config *Config) (err error) {
	l, err := config.List()
	//fmt.Println("all list:", l)
	if err != nil {
		return
	}

	for _, s := range l {
		err = os.Remove(s)
		if err != nil {
			return
		}
	}
	return nil
}
