package main

import(
	"github.com/kpiotrowski/go_watchdog/watchdog"
	"log"
)

func main() {

	//TODO ADD CONFIGURABLE VARIABLES AND RUN AS DEMON
	service, err:= watchdog.NewService("mysql", "10s", "10s", 4)
	if err != nil{
		log.Fatal(err)
	}
	service.Watch()

}