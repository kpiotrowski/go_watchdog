package mail

import (
	"github.com/kpiotrowski/go_watchdog/common"
	"net/mail"
)

type SenderInterface interface {
	Send(title string, body []byte) error
}

type emailSender struct {

}

func NewEmailSender(config common.MailConf) (*emailSender, error) {
	_, _ = mail.ParseAddress(config.MailFromAddress)


	return nil, nil
}

func (e *emailSender) Send(title string, body []byte) error {

	return nil
}