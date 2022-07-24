package mail

import (
	"os"
	"testing"
)

// Testing for this one is a little more difficult, since SendMessage has to send out an email.
// You need SES credentials for the test to run, as well as a verified email to send from.
// Replace testHost data with your host/sender information, and testRecipient with a known resolvable
// target email.

var testHost Host = Host{
	Username: os.Getenv("SES_USERNAME"),
	Password: os.Getenv("SES_PASSWORD"),
	Host:     "email-smtp.us-east-2.amazonaws.com",
	Port:     587,
	Sender:   "jnichols2719@protonmail.com",
}

var testRecipient string = "success@simulator.amazonses.com"

func TestSendMessage(t *testing.T) {
	msg := NewAuthMessage(testRecipient, "AAAAAA")
	SendMessage(testHost, testRecipient, msg)
}
