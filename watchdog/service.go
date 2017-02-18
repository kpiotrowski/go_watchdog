package watchdog

import (
	"time"
	"fmt"
	"errors"
	"os"
)

type service interface {
	running()
	start()
	Watch()
}

type serviceStruct struct {
	name string
	checkInterval time.Duration
	restartInterval time.Duration
	restartNumber int
}


func NewService(name, checkInterval, startInterval string, tries int) (*serviceStruct, error) {
	if _, err := os.Stat("/etc/init.d/"+name); err != nil {
		return nil, errors.New(fmt.Sprintf("Service %s doesn't exist\n",name))
	}

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

	service := new(serviceStruct)
	service.name = name
	service.checkInterval = checkI
	service.restartInterval = restartI
	service.restartNumber = tries

	return service, nil
}

func (service *serviceStruct) running() {

}

func (service *serviceStruct) start() {

}

func (service *serviceStruct) Watch() {

}