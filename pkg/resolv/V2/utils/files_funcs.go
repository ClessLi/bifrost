package utils

import "fmt"

func RemoveFiles(files []string) error {
	for _, path := range files {
		/*err := os.Remove(path)
		if err != nil {
			return err
		}*/
		// debug test
		fmt.Printf("remove: %s", path)
		// debug test end
	}
	return nil
}
