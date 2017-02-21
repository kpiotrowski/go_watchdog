# go_watchdog

![Build Status](https://travis-ci.org/kpiotrowski/go_watchdog.svg?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/kpiotrowski/go_watchdog/badge.svg?branch=master)](https://coveralls.io/github/kpiotrowski/go_watchdog?branch=master)

Simple watchdog written in golang

## RUN
Program is running as a daemon, and works if user is not logged.
You can start multiple go_watchdogs for different services.
 
To disable go_watchdog you need to send SIGQUIT or SIGTERM signal to the process. Go_watchdog catches signal and turns off.
Go_watchdog also turns off if service was not started after number of attempts. 


## USAGE
```
./go_watchdog
-a int
    Number of attempts to start service (default 4)
-c string
    Service status check interval [duration string] (default "60s")
-i string
    Service start interval [duration string] (default "10s")
-l string
    Log fle name (default "log")
-m string
    File name with mail config (default "mail.conf")
-s string
    service name to watch

```

## NOTIFICATIONS

Go_watchdog send emails and log results for given events:

- service is down
- service started after number of attempts
- service cannot be started after number of attempts

For sending emails you need to create mail config file:
```
[Mail]
MailFromAddress = "from@gmail.com"
MailFromPassword = "pass"
MailServerAddress = "smtp.gmail.com:587"
MailTo = "to@gmail.com"
```

## EXAMPLES

Watch if mysql is running every 45s. If not try to start every 10s up to 5 times:

`./go_watchdog -s mysql -c 45s -i 10s -a 5`

Watch docker and change default config and log file:

`./go_watchdog -s docker -l docker_log -m mailConfig`
