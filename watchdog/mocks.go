package watchdog

import (
	"github.com/stretchr/testify/mock"
	"os"
)

type mockedOs struct {
	mock.Mock
}

type mockedSender struct {
	mock.Mock
}

func (m *mockedOs) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	return nil, args.Error(1)
}

func (m *mockedOs) ExecOutput(name string, arg ...string) ([]byte, error) {
	args := m.Mock.Called(name, arg)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockedSender) 	Send(title string, body []byte) error{
	args := m.Mock.Called()
	return args.Error(0)
}

