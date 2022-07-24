package mail

import (
	"fmt"
	"net/smtp"
)

// struct Host holds host data for sending SMTP through SES.
type Host struct {
	Username string
	Password string
	Host     string
	Port     int
	Sender   string
}

// Generate a PlainAuth to use with smtp.SendMail using host info
func (h Host) plainAuth() smtp.Auth {
	return smtp.PlainAuth("", h.Username, h.Password, h.Host)
}

// Generate an address for smtp.SendMail
func (h Host) address() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

// Generate a new authentication message given a target email and authentication code.
func NewAuthMessage(sendTo string, authCode string) []byte {
	msg := fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: Authentication Code\r\n"+
			"\r\n"+
			"Your authentication code is %s.\n"+
			"This code will expire in 1 minute.\r\n",
		sendTo, authCode,
	)
	return []byte(msg)
}

// Use smtp to send a message through the target SES host
func SendMessage(sendFrom Host, sendTo string, msg []byte) error {
	auth := sendFrom.plainAuth()
	target := []string{sendTo}
	addr := sendFrom.address()
	return smtp.SendMail(addr, auth, sendFrom.Sender, target, msg)
}
