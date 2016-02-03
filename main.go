package main

import (
	"fmt"
	"github.com/zero-boilerplate/go-api-helpers/service"
	"path/filepath"

	service2 "github.com/ayufan/golang-kardianos-service"
)

type app struct {
	logger             service2.Logger
	watcherDoneChannel chan bool
}

func (a *app) OnStop() {
	defer recover()
	close(a.watcherDoneChannel)
}

func (a *app) Run(logger service2.Logger) {
	a.logger = logger
	defer func() {
		if r := recover(); r != nil {
			a.logger.Errorf("Run app error: %s", getStringFromRecovery(r))
		}
	}()

	userHomeDir := getUserHomeDir()
	watchDir := filepath.Join(userHomeDir, ".script-watcher", "scripts")
	if !doesDirExist(watchDir) {
		panic(fmt.Sprintf("The watch dir '%s' does not exist", watchDir))
		return
	}

	a.scanDirForExistingFile(watchDir)
	a.startWatching(watchDir)
}

func main() {
	a := &app{}
	service.NewServiceRunnerBuilder("Script Watcher", a).WithOnStopHandler(a).WithServiceUserName_AsCurrentUser().Run()
}
