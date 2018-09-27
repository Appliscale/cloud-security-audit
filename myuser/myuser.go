package myuser

import (
	"os/user"
)

// GetUserHomeDir gets path to home directory.
func GetUserHomeDir() (string, error) {
	myUser, userError := user.Current()
	if userError != nil {
		return "", userError
	}
	path := myUser.HomeDir

	return path, nil
}
