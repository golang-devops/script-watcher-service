package main

import (
	"os/exec"
	"strings"
	"time"
)

type executor struct {
	exe  string
	args []string
}

func (e *executor) generateExeAndArgs(scriptFile string) (outExe string, outArgs []string) {
	if e.exe == "" {
		return scriptFile, []string{}
	} else {
		outExe = e.exe
		outArgs = e.args
		outArgs = append(outArgs, scriptFile)
		// outArgs = append(outArgs, args...)
		return
	}
}

func (e *executor) execute(a *app, scriptFile string) {
	defer a.deleteFile(scriptFile)

	sleepDur := 2 * time.Second
	a.logger.Infof("Sleeping for '%s' to allow file disk-write to complete for file '%s'", sleepDur.String(), scriptFile)
	time.Sleep(sleepDur)

	exe, args := e.generateExeAndArgs(scriptFile)
	cmd := exec.Command(exe, args...)

	out, err := cmd.CombinedOutput()
	cleanedOutput := strings.Replace(strings.Replace(string(out), "\n", "\\n", -1), "\r", "", -1)
	if err != nil {
		a.logger.Errorf("Could not execute script '%s'. ERROR: '%s', CombinedOutput: '%s'", scriptFile, err.Error(), cleanedOutput)
		return
	}

	a.logger.Infof("Output of script '%s' was '%s'", scriptFile, cleanedOutput)
}

func getScriptExecutor(scriptFilePath string) *executor {
	lowerPath := strings.ToLower(scriptFilePath)

	if strings.HasSuffix(lowerPath, ".bat") {
		return &executor{"cmd", []string{"/c"}}
	} else if strings.HasSuffix(lowerPath, ".sh") {
		return &executor{"bash", []string{"-c"}}
	} else if strings.HasSuffix(lowerPath, ".py") {
		return &executor{"python", []string{}}
	} else if strings.HasSuffix(lowerPath, ".go") {
		return &executor{"go", []string{"run"}}
	}

	return nil
}
