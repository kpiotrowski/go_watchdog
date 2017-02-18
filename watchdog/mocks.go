package watchdog

import (
	"github.com/stretchr/testify/mock"
	"os"
)

type mockedOs struct {
	mock.Mock
}

func (m *mockedOs) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(os.FileInfo), args.Error(1)
}
func (m*mockedOs) ExecOutput(name string, arg ...string) ([]byte, error) {
	args := m.Called(name, arg)
	return args.Get(0).([]byte), args.Error(1)
}
