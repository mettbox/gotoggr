package util

import "os"

func FileExists(path string) (found bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		found = true
	}

	return
}
