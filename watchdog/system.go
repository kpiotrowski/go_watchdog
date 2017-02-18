package watchdog

/**
	Needed to mock file system for test purposes
 */

import (
	"os"
	"os/exec"
)

type osInterface interface {
	Stat(name string) (os.FileInfo, error)
	ExecOutput(name string, arg ...string) ([]byte, error)
}
type OsFS struct {}

func (OsFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OsFS) ExecOutput(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}
