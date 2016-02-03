package main

import (
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
)

func (a *app) handleFile(filePath string) {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		a.logger.Errorf("Cannot get ABS path of script file '%s', error: '%s'", filePath, err.Error())
		return
	}

	executor := getScriptExecutor(absPath)

	if executor == nil {
		a.logger.Warningf("Non-script file '%s' was found", absPath)
		return
	}

	go executor.execute(a, absPath)
}

func (a *app) scanDirForExistingFile(watchDir string) {
	e := filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		a.handleFile(path)

		return nil
	})
	if e != nil {
		a.logger.Warningf("Unable to walk dir '%s', error: '%s'", watchDir, e.Error())
	}
}

func (a *app) startWatching(watchDir string) {
	watcher, err := fsnotify.NewWatcher()
	checkError(err)
	defer watcher.Close()

	a.watcherDoneChannel = make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					a.handleFile(ev.Name)
				} else if !ev.IsDelete() {
					a.logger.Warningf("Non create/delete event: %s", ev.String())
				}
				break
			case e := <-watcher.Error:
				a.logger.Errorf("Watcher error: %s", e.Error())
			}
		}
	}()

	a.logger.Infof("Now watching dir '%s'", watchDir)
	err = watcher.Watch(watchDir)
	checkError(err)

	<-a.watcherDoneChannel
}
