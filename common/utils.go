package common

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	logger "github.com/Sirupsen/logrus"
)

func Exec(command string, args ...string) (output string, err error) {
	var commandArgs []string
	var spaceRegexp = regexp.MustCompile("[\\s]+")

	commands := spaceRegexp.Split(command, -1)
	command = commands[0]

	if len(commands) > 1 {
		commandArgs = commands[1:]
	}

	if len(args) > 0 {
		commandArgs = append(commandArgs, args...)
	}

	fullCommand, err := exec.LookPath(command)

	if err != nil {
		return "", fmt.Errorf("%s cannot be found", command)
	}

	cmd := exec.Command(fullCommand, commandArgs...)
	cmd.Env = os.Environ()

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	out, err := cmd.Output()

	if err != nil {
		logger.Debug(fullCommand, " ", strings.Join(commandArgs, " "))
		err = errors.New(stdErr.String())
		return
	}

	output = strings.Trim(string(out), "\n")
	return
}

func CreateDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}
}
