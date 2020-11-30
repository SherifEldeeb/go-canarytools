// +build !windows

package canarytools

import (
	"fmt"
)

func SetFileAttributeHiddenAndSystem(filename string) error {
	fmt.Println("HideFile works only on windows builds")
	return nil
}

func SetFileAttributeSystem(filename string) error {
	fmt.Println("HideFile works only on windows builds")
	return nil
}
