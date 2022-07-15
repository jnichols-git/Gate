package main

import (
	"auth/pkg/ses_mail"
	"os"
)

func main() {
	host := ses_mail.Host{
		Username: os.Getenv("SES_USERNAME"),
		Password: os.Getenv("SES_PASSWORD"),
		Hostname: os.Getenv("SES_ENDPOINT"),
	}
	msg := ses_mail.NewAuthMessage("jani9652@colorado.edu", "AAAAAA")
	ses_mail.SendMessage(host, "jani9652@colorado.edu", msg)
}
