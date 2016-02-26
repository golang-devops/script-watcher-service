package main

import (
	"sync"
	"time"
)

type fileEventManager struct {
	sync.RWMutex
	fileEventTimes map[string]time.Time
}

func (f *fileEventManager) cleanupFileEventTimes() {
	listToRemove := []string{}
	for fileName, eventTime := range f.fileEventTimes {
		if time.Now().Sub(eventTime) > 1*time.Second {
			listToRemove = append(listToRemove, fileName)
		}
	}

	for _, fileName := range listToRemove {
		delete(f.fileEventTimes, fileName)
	}
}

func (f *fileEventManager) isDuplicateEvent(fileName string) bool {
	f.Lock()
	defer f.Unlock()

	f.cleanupFileEventTimes()

	_, ok := f.fileEventTimes[fileName]

	f.fileEventTimes[fileName] = time.Now()
	return !ok
}
