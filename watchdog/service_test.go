package watchdog

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"errors"
	"os"
	"os/exec"
)

const (
	incorrectName = "name"
	correctName = "mysql"
	incorrectInterval = "interval"
	correctInterval = "1ms"
	incorrectTries = 0
	correctTries = 10
)

func TestSystemWrapperStat(t *testing.T){
	osFS := OsFS{}
	expected, err := os.Stat("log")
	got, errGot := osFS.Stat("log")
	assert.Equal(t, expected, got, "Results are different")
	assert.Equal(t, err, errGot, "Results are different")
}

func TestSystemWrapperExec(t *testing.T){
	osFS := OsFS{}
	expected, err := exec.Command("ls").Output()
	got, errGot := osFS.ExecOutput("ls")
	assert.Equal(t, expected, got, "Results are different")
	assert.Equal(t, err, errGot, "Results are different")
}

func TestNewServiceIncorrectCheckInterval(t *testing.T) {
	_, err := NewService(correctName, incorrectInterval, correctInterval, correctTries)
	assert.NotNil(t, err, "Incorrect check interval should return error")
}

func TestNewServiceIncorrectStartInterval(t *testing.T) {
	_, err := NewService(correctName, correctInterval, incorrectInterval, correctTries)
	assert.NotNil(t, err, "Incorrect start interval should return error")
}

func TestNewServiceIncorrectTriesNumber(t *testing.T) {
	_, err := NewService(correctName, correctInterval, correctInterval, incorrectTries)
	assert.NotNil(t, err, "Incorrect stries number should return error")
}

func TestNewServiceIncorrectServiceName(t *testing.T) {
	osMock := new(mockedOs)
	osMock.On("Stat", initDPath+incorrectName).Return(nil, errors.New("incorect service name"))
	_, err := newServiceWithOs(incorrectName, correctInterval, correctInterval, correctTries, osMock)
	assert.NotNil(t, err, "Incorrect name should return error")
}

func TestNewServiceEmptyServiceName(t *testing.T) {
	osMock := new(mockedOs)
	osMock.On("Stat", initDPath).Return(nil, errors.New("empty service name"))
	_, err := newServiceWithOs("", correctInterval, correctInterval, correctTries, osMock)
	assert.NotNil(t, err, "Incorrect name should return error")
}

func TestNewServiceSuccess(t *testing.T) {
	osMock := new(mockedOs)
	osMock.On("Stat", initDPath+correctName).Return(nil, nil)
	service, err := newServiceWithOs(correctName, correctInterval, correctInterval, correctTries, osMock)
	assert.Nil(t, err, "Shouldn't return erorr")
	assert.NotNil(t, service, "Sould return new service")
	assert.Equal(t,correctName, service.name, "Service name is incorrect")
	interval , _ := time.ParseDuration(correctInterval)
	assert.Equal(t, interval, service.checkInterval, "Service check interval is incorrect")
	assert.Equal(t, interval, service.startInterval, "Service start interval is incorrect")
	assert.Equal(t, correctTries, service.startTries, "Service tries is incorrect")
}

func newTestService() serviceStruct {
	interval , _ := time.ParseDuration(correctInterval)
	s := serviceStruct{name:correctName,checkInterval:interval,startTries:correctTries,startInterval:interval}
	return s
}

func TestServiceStruct_RunningFalse(t *testing.T) {
	testService := newTestService()
	osMock := new(mockedOs)
	osMock.On("ExecOutput", serviceCommand, []string{testService.name, statusCommand}).Return([]byte{}, errors.New("error"))
	testService.os = osMock
	running := testService.Running()
	assert.False(t, running, "Service running should return false")
}

func TestServiceStruct_RunningTrue(t *testing.T) {
	testService := newTestService()
	osMock := new(mockedOs)
	osMock.On("ExecOutput", serviceCommand, []string{testService.name, statusCommand}).Return([]byte{}, nil)
	testService.os = osMock
	running := testService.Running()
	assert.True(t, running, "Service running should return true")
}

func TestServiceStruct_StartFalse(t *testing.T) {
	testService := newTestService()
	osMock := new(mockedOs)
	osMock.On("ExecOutput", serviceCommand, []string{testService.name, startCommand}).Return([]byte{}, errors.New("error"))
	testService.os = osMock
	running := testService.Start()
	assert.False(t, running, "Service start should return false")
}

func TestServiceStruct_StartTrue(t *testing.T) {
	testService := newTestService()
	osMock := new(mockedOs)
	osMock.On("ExecOutput", serviceCommand, []string{testService.name, startCommand}).Return([]byte{}, nil)
	testService.os = osMock
	running := testService.Start()
	assert.True(t, running, "Service start should return true")
}

func TestServiceStruct_WatchBreak(t *testing.T) {
	testService := newTestService()
	osMock := new(mockedOs)
	osMock.On("ExecOutput", serviceCommand, []string{testService.name, statusCommand}).Return([]byte{}, nil)
	testService.os = osMock
	mockedSender := mockedSender{}
	stop := make(chan bool)
	done := make(chan bool)
	go func(){
		err := testService.Watch(&mockedSender, stop)
		assert.Nil(t,err,"Watch error should be nil")
		done <- true
	}()
	stop <- true
	<- done
}

func TestServiceStruct_WatchStopAfterFailedStarts(t *testing.T) {
	testService := newTestService()
	osMock := new(mockedOs)
	osMock.On("ExecOutput", serviceCommand, []string{testService.name, statusCommand}).Return([]byte{}, errors.New("error"))
	osMock.On("ExecOutput", serviceCommand, []string{testService.name, startCommand}).Return([]byte{},  errors.New("error"))
	testService.os = osMock
	mockedSender := mockedSender{}
	mockedSender.On("Send").Return(errors.New(""))
	stop := make(chan bool)
	err := testService.Watch(&mockedSender, stop)
	assert.NotNil(t,err,"Multiple failed starts should stop watchdog")
}
