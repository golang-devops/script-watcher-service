package main

import (
	"fmt"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getStringFromRecovery(r interface{}) string {
	switch t := r.(type) {
	case error:
		return t.Error()
	}
	return fmt.Sprintf("%#v", r)
}
