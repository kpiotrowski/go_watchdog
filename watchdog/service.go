package watchdog

import (
	"time"
	"fmt"
	"errors"
)

const (
	initDPath = "/etc/init.d/"
)

type service interface {
	running()
	start()
	Watch()
}

type serviceStruct struct {
	name string
	checkInterval time.Duration
	startInterval time.Duration
	startTries int
	os osInterface
}

func NewService(name, checkInterval, startInterval string, tries int) (*serviceStruct, error) {
	return newServiceWithOs(name, checkInterval, startInterval, tries, OsFS{})
}

func newServiceWithOs(name, checkInterval, startInterval string, tries int, os osInterface) (*serviceStruct, error) {
	checkI, err := time.ParseDuration(checkInterval)
	if err != nil{
		return  nil, err
	}

	restartI, err := time.ParseDuration(startInterval)
	if err != nil{
		return  nil, err
	}

	if tries < 1 {
		return nil, errors.New("Incorrect number of tries to run service")
	}

	if _, err := os.Stat(initDPath+name); err != nil {
		return nil, errors.New(fmt.Sprintf("Service %s doesn't exist\n",name))
	}

	service := new(serviceStruct)
	service.name = name
	service.checkInterval = checkI
	service.startInterval = restartI
	service.startTries = tries
	service.os = os

	return service, nil
}

func (service *serviceStruct) running() {

}

func (service *serviceStruct) start() {

}

func (service *serviceStruct) Watch() {

}