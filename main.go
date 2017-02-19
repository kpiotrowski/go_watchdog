package main

import(
	"log"
	"github.com/kpiotrowski/go_watchdog/common"
	"github.com/kpiotrowski/go_watchdog/mail"
	"github.com/kpiotrowski/go_watchdog/watchdog"
	"flag"
	"github.com/sevlyar/go-daemon"
)

var(
	signal = flag.String("k", "", `send signal to the daemon
                quit — graceful shutdown
                stop — fast shutdown`)
	processPid = flag.Int("p", 0, `process pid used to send signal`)
	serviceName = flag.String("s", "", "service name to watch")
	logFile = flag.String("l","log", "log fle name")
	checkInterval = flag.String("c", "60s", "Check isterval [duration string]")
	startInterval = flag.String("i", "10s", "Start interval [duration string]")
	attempts = flag.Int("a", 4, "Number of attempts when starting service")
	mailFile = flag.String("m", "mail.conf", "File name with mail config")
)

func main() {
	flag.Parse()

	conf, err := common.LoadConfig(*mailFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	mailSender, err := mail.NewEmailSender(conf.Mail)
	if err != nil {
		log.Fatal(err)
		return
	}
	service, err:= watchdog.NewService(*serviceName, *checkInterval, *startInterval, *attempts)
	if err != nil{
		log.Fatal(err)
	}

	context := daemon.Context{
		LogFileName: *logFile,
		LogFilePerm: 0640,
	}

	child, err := context.Reborn()
	if err != nil {
	      log.Fatal(err)
	}
	if child != nil {     //nil is returned in child process
		return
	}
	defer context.Release()

	go service.Watch(mailSender)

	err = daemon.ServeSignals()
	if err != nil {
	      log.Println("Serve siglans error:", err)
	}

	service.Watch(mailSender)
}