package environment

import (
	"os/user"
)

func GetUserHomeDir() (string, error) {
	myUser, userError := user.Current()
	if userError != nil {
		return "", userError
	}
	path := myUser.HomeDir

	return path, nil
}

func UniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}
	return us
}
