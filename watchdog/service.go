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
	serviceDown = iota
	serviceStart = iota
	serviceCannotStart = iota
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

	if len(name) < 1 {
		return nil, errors.New("Service name cannot be empty")
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

func notify(sender mail.SenderInterface, service string, attempts, status int){
	var logMsg string
	var title string

	switch status {
	case serviceDown:
		logMsg = fmt.Sprintf("%s Service %s is down", time.Now().String(), service)
		title = service+" is down"
	case serviceStart:
		logMsg = fmt.Sprintf("%s Service %s started after %d attempts",time.Now().String(), service, attempts)
		title = service+" started"
	case serviceCannotStart:
		logMsg = fmt.Sprintf("%s Service %s can't be started after %d attempts", time.Now().String(), service, attempts)
		title = service+" start failed"
	}

	log.Println(logMsg)
	go sender.Send(title, []byte(logMsg))
}

func (service *serviceStruct) Watch(sender mail.SenderInterface, stopChan chan bool) {
	loop := true
	checInterval := make(chan time.Time)
	startInterval := make(chan time.Time)
	go func() {
		<- stopChan
		loop = false
		checInterval <- time.Now()
		startInterval <- time.Now()
	}()

	for loop {
		run := service.Running()
		if !run {
			notify(sender, service.name, 0, serviceDown)
			for i:=1 ; i<= service.startTries && loop ; i++ {
				if run = service.Start() ; run {
					notify(sender, service.name, i, serviceStart)
					break
				}
				go func(){
					time.Sleep(service.startInterval)
					startInterval <- time.Now()
				}()
				<-startInterval
			}
			if !run {
				notify(sender, service.name, service.startTries, serviceCannotStart)
				return
			}
		}
		if !loop {
			return
		}
		go func(){
			time.Sleep(service.checkInterval)
			checInterval <- time.Now()
		}()
		<-checInterval
	}
}