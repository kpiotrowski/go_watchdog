package watchdog

import (
	"time"
	"fmt"
	"errors"
	"log"
	"github.com/kpiotrowski/go_watchdog/mail"
)

const (
	initDPath = "/etc/init.d/"
	serviceCommand = "service"
	startCommand = "start"
	statusCommand = "status"
)

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

func (service *serviceStruct) Watch(sender mail.SenderInterface) {
	for {
		run := service.Running()
		if !run {
			logMsg := fmt.Sprintf("%s Service %s is down", time.Now().String(), service.name)
			log.Println(logMsg)
			sender.Send(service.name+" is down", []byte(logMsg))

			for i:=1 ; i<= service.startTries ; i++ {
				if run = service.Start() ; run {
					logMsg = fmt.Sprintf("%s Service %s started after %d attempts",time.Now().String(), service.name, i)
					log.Println(logMsg)
					sender.Send(service.name+" started",[]byte(logMsg))
					break
				}
				time.Sleep(service.startInterval)
			}
			if !run {
				logMsg = fmt.Sprintf("%s Service %s can't be started after %d attempts", time.Now().String(), service.name, service.startTries)
				log.Println(logMsg)
				sender.Send(service.name+" start Failed", []byte(logMsg))
				return
			}

		}
		time.Sleep(service.checkInterval)
	}
}