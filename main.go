package main

import(
	"log"
	"github.com/kpiotrowski/go_watchdog/common"
	"github.com/kpiotrowski/go_watchdog/mail"
	"fmt"
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
	err = mailSender.Send("tit",[]byte("body"))
	if err != nil {
		fmt.Print(err)
	}

	////TODO ADD CONFIGURABLE VARIABLES AND RUN AS DEMON
	//service, err:= watchdog.NewService("mysql", "10s", "10s", 4)
	//if err != nil{
	//	log.Fatal(err)
	//}
	//service.Watch(mailSender)
}