package watchdog

import (
	"time"
	"fmt"
	"errors"
	"log"
)

const (
	initDPath = "/etc/init.d/"
	serviceCommand = "service"
	startCommand = "start"
	statusCommand = "status"
)

type service interface {
	Running() bool
	Start() (bool, error)
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

func (service *serviceStruct) Running() bool {
	_, err := service.os.ExecOutput(serviceCommand, service.name, statusCommand)
	if err != nil {
		return false
	}
	return true
}

func (service *serviceStruct) Start() bool {
	_, err := service.os.ExecOutput(serviceCommand, service.name, startCommand)
	if err != nil {
		return false
	}
	return true
}

func (service *serviceStruct) Watch() {
	for {
		run := service.Running()
		if !run {
			log.Println(fmt.Sprintf("%s Service %s is down", time.Now().String(), service.name))
		}
		for !run {
			for i:=1 ; i<= service.startTries && !run ; i++ {
				run = service.Start()
				if run {
					log.Println(fmt.Sprintf("%s Service %s started after %d attempts",time.Now().String(), service.name, i))
					break
				}
			}
			if !run {
				log.Println(fmt.Sprintf("%s Service %s can't be started after %d attempts", time.Now().String(), service.name, service.startTries))
				time.Sleep(service.startInterval)
			}
		}
		time.Sleep(service.checkInterval)
	}
}