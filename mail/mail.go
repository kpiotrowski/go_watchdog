package mail

import (
	"github.com/kpiotrowski/go_watchdog/common"
	"net/mail"
	"net/smtp"
	"net"
	"errors"
	"fmt"
)

type SenderInterface interface {
	Send(title string, body []byte) error
}

type emailSender struct {
	from *mail.Address
	smtpServerPort string
	smtpServerHost string
	to *mail.Address

	auth smtp.Auth
}

func NewEmailSender(config common.MailConf) (*emailSender, error) {
	sender := new(emailSender)
	addrFrom, err := mail.ParseAddress(config.MailFromAddress)
	if err != nil {
		return nil, err
	}
	sender.from = addrFrom
	addrTo, err := mail.ParseAddress(config.MailTo)
	if err != nil {
		return nil, err
	}
	sender.to = addrTo
	host, port, err := net.SplitHostPort(config.MailServer)
	if err != nil {
		return nil, err
	}
	sender.smtpServerHost = host
	sender.smtpServerPort = port
	if len(config.MailFromAddress) < 1 {
		return nil, errors.New("Email password cannot be empty")
	}
	sender.auth = smtp.PlainAuth("", addrFrom.Address, config.MailFromPassword, host)

	return sender, nil
}

func createMessage(title string, body []byte) []byte{
	byteMsg := []byte(fmt.Sprintf("Subject: %s\r\n",title))
	return append(byteMsg, body...)
}

func (e *emailSender) Send(title string, body []byte) error {
	err := smtp.SendMail(
		e.smtpServerHost+":"+e.smtpServerPort,
		e.auth,
		e.from.Address,
		[]string{e.to.Address},
		createMessage(title,body),
	)
	return err
}