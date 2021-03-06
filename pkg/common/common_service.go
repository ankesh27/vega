package common

import (
	"fmt"
	"os"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func EnsureDir(dir string) error {
	if info, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("Could not create %s: %s", dir, err)
		}
	} else if !info.IsDir() {
		// CHECK: Why not errors.New()?
		return fmt.Errorf("%s must be a directory", dir)
	}

	return nil
}
