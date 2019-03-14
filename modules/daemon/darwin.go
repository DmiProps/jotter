package daemon

import (
	// System packages.
	"os/exec"
	"regexp"
	"strconv"

	// Jotter packages.
	"github.com/dmiprops/jotter/modules/setting"
)

// StartDaemon starts jotter service.
func StartDaemon(args ...string) error {
	newargs := []string{
		"submit",
		"-l", setting.DemonLabel,
		"-o", setting.RunDir+"/log/"+setting.DaeminLogFile,
		"-e", setting.RunDir+"/log/"+setting.DaeminErrFile,
		"--", setting.RunDir+"/jotter",
	}

	return exec.Command("launchctl", append(newargs, args ...) ...).Start()
}

// StopDaemon stops jotter service.
func StopDaemon() error {
	return exec.Command(
		"launchctl",
		"remove", setting.DemonLabel,
	).Start()
}

// RestartDaemon restarts jotter service.
func RestartDaemon(args ...string) error {
	err := StopDaemon()
	if err != nil {
		return err
	}
	
	return StartDaemon(args ...)
}

// GetPID returns PID jotter service.
func GetPID() (int, error) {
	output, err := exec.Command(
		"launchctl",
		"list", setting.DemonLabel,
	).Output()
	if err != nil {
		return 0, err
	}

	reg := regexp.MustCompile(`("PID" = ).+(;)`)
	str := reg.FindString(string(output))
	if str == "" {
		return 0, nil
	}

	return strconv.Atoi(str[8:len(str)-1])
}

// CheckRun returns true if jotter service is run.
func CheckRun() bool {
	PID, _ := GetPID()
	return PID != 0
}