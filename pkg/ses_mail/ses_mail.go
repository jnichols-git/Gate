package ses_mail

import (
	"fmt"
	"net/smtp"
)

// struct Host holds host location/auth information for smtp
type Host struct {
	Username string // Host username
	Password string // Host password
	Hostname string // Host name/address (ex. smtp.gmail.com)
}

func (h Host) PlainAuth() smtp.Auth {
	return smtp.PlainAuth("", h.Username, h.Password, h.Hostname)
}

func (h Host) Address() string {
	return fmt.Sprintf("%s:%d", h.Hostname, 587)
}

func NewAuthMessage(sendTo, authCode string) []byte {
	msg := fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: Authentication Code\r\n"+
			"\r\n"+
			"Your authentication code is %s\r\n",
		sendTo, authCode,
	)
	return []byte(msg)
}

func SendMessage(sendFrom Host, sendTo string, msg []byte) error {
	auth := sendFrom.PlainAuth()
	target := []string{sendTo}
	addr := sendFrom.Address()
	err := smtp.SendMail(addr, auth, "jnichols2719@protonmail.com", target, msg)
	if err != nil {
		return err
	}
	return nil
}
