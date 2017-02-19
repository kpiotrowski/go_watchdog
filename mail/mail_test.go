package mail

import (
	"github.com/kpiotrowski/go_watchdog/common"
	"testing"
	"github.com/stretchr/testify/assert"
	"net/mail"
	"net"
	"net/smtp"
)

const (
	msgTitle  = "title"
	msgBody = "body"
	emailAddrCorrect = "name@server.com"
	emailAddrIncorrect = "incorrec"
	serverCorrect = "smtp.gmail.com:578"
	serverIncorrect = "serverIncorrect"
	passwordCorrect = "password"
)

func TestCreateEmailMessage (t *testing.T){
	expected := []byte("Subject: "+msgTitle+"\n\n"+msgBody)
	msg := createMessage(msgTitle, []byte(msgBody))
	assert.Equal(t, msg, expected, "Email msg is incorrect")
}

func newTestConfig () common.MailConf {
	return common.MailConf{
		MailFromAddress:emailAddrCorrect,
		MailFromPassword:passwordCorrect,
		MailServer: serverCorrect,
		MailTo:emailAddrCorrect,
	}
}

func TestNewEmailSenderIncorrectFromAddr(t *testing.T) {
	testConf := newTestConfig()
	testConf.MailFromAddress = emailAddrIncorrect
	_, err := NewEmailSender(testConf)
	assert.NotNil(t,err, "incorrect From address should return error")
}


func TestNewEmailSenderIncorrectToAddr(t *testing.T) {
	testConf := newTestConfig()
	testConf.MailTo = emailAddrIncorrect
	_, err := NewEmailSender(testConf)
	assert.NotNil(t, err, "incorrect To address should return error")
}


func TestNewEmailSenderIncorrectPassword(t *testing.T) {
	testConf := newTestConfig()
	testConf.MailFromPassword = ""
	_, err := NewEmailSender(testConf)
	assert.NotNil(t, err, "incorrect password should return error")
}

func TestNewEmailSenderIncorrectServer(t *testing.T) {
	testConf := newTestConfig()
	testConf.MailServer = serverIncorrect
	_, err := NewEmailSender(testConf)
	assert.NotNil(t, err, "incorrect server address should return error")
}

func TestNewEmailSenderSuccess(t *testing.T) {
	sender, err := NewEmailSender(newTestConfig())
	assert.Nil(t, err, "Shouldn't return error for correct config")
	assert.NotNil(t, sender, "Sender shouldn't be null for correct config")
	address, _ := mail.ParseAddress(emailAddrCorrect)
	host, port, _ := net.SplitHostPort(serverCorrect)
	auth := smtp.PlainAuth("",address.Address, passwordCorrect, host)
	assert.Equal(t, sender.from, address, "From address is incorrect")
	assert.Equal(t, sender.to, address, "To address is incorrect")
	assert.Equal(t, sender.smtpServerHost, host, "Host is incorrect")
	assert.Equal(t, sender.smtpServerPort, port, "Port is incorrect")
	assert.Equal(t, sender.auth, auth, "Auth is incorrect")
}