package canarytools

import (
	"fmt"
	"os"
)

func fileExists(filename string) (exists bool, err error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if info != nil {
		return !info.IsDir(), nil
	}
	return false, fmt.Errorf("os.Stat returned a nil 'FileInfo' struct")
}
