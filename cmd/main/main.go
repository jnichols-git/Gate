package main

import (
	"auth/pkg/authcode"
	"auth/pkg/ses_mail"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	host := ses_mail.Host{
		Username: os.Getenv("SES_USERNAME"),
		Password: os.Getenv("SES_PASSWORD"),
		Host:     "email-smtp.us-east-2.amazonaws.com",
		Port:     587,
		Sender:   "jnichols2719@protonmail.com",
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Please enter email to authenticate: ")
	email, _ := reader.ReadString('\n')
	code := authcode.NewAuthCode(email)
	msg := ses_mail.NewAuthMessage("jani9652@colorado.edu", code)
	ses_mail.SendMessage(host, "jani9652@colorado.edu", msg)
	fmt.Printf("Please enter authentication code: ")
	c, _ := reader.ReadString('\n')
	c = strings.Trim(c, "\r\n")
	valid := authcode.ValidateAuthCode(email, c)
	if valid {
		fmt.Println("Valid!")
	} else {
		fmt.Println("Not valid.")
	}
}
