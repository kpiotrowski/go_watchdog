package main

import(
	"github.com/kpiotrowski/go_watchdog/watchdog"
	"log"
	"github.com/kpiotrowski/go_watchdog/common"
	"github.com/kpiotrowski/go_watchdog/mail"
)

func main() {
	conf, err := common.LoadConfig("mail.conf")
	if err != nil {
		log.Fatal(err)
		return
	}
	mailSender, err := mail.NewEmailSender(conf.Mail)
	if err != nil {
		log.Fatal(err)
		return
	}

	//TODO ADD CONFIGURABLE VARIABLES AND RUN AS DEMON
	service, err:= watchdog.NewService("mysql", "10s", "10s", 4)
	if err != nil{
		log.Fatal(err)
	}
	service.Watch(mailSender)
}