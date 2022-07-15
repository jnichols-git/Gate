package main

import (
	"auth/pkg/ses_mail"
	"os"
)

func main() {
	host := ses_mail.Host{
		Username: os.Getenv("SES_USERNAME"),
		Password: os.Getenv("SES_PASSWORD"),
		Host:     "email-smtp.us-east-2.amazonaws.com",
		Port:     587,
		Sender:   "jnichols2719@protonmail.com",
	}
	msg := ses_mail.NewAuthMessage("jani9652@colorado.edu", "AAAAAA")
	ses_mail.SendMessage(host, "jani9652@colorado.edu", msg)
}
