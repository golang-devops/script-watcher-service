package main

import (
	"os"
	"os/user"
)

func getUserHomeDir() string {
	u, err := user.Current()
	if err != nil {
		panic("Cannot get current user: " + err.Error())
	}

	return u.HomeDir
}

func doesDirExist(dir string) bool {
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic("Unexpected dir stat error for dir '" + dir + "', error: " + err.Error())
	}
	return true
}

func (a *app) deleteFile(p string) {
	err := os.Remove(p)
	if err != nil {
		a.logger.Warningf("Unable to remove file '%s', error: %s", p, err.Error())
	}
}
