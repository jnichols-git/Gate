package main

import (
	"auth/pkg/authcode"
	"auth/pkg/authjwt"
	"auth/pkg/ses_mail"
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var secret []byte = []byte("test secret")

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
	fmt.Printf("Please enter authentication token: ")
	token, _ := reader.ReadString('\n')
	if strings.Trim(token, "\r\n") != "" {
		jwt := authjwt.UnmarshalJWT(token)
		valid := authjwt.VerifySignature(jwt, secret)
		if valid {
			fmt.Printf("Authentication valid. Welcome, %s\n", jwt.Body.ForUser)
			return
		}
	}
	fmt.Printf("Please enter email to authenticate: ")
	email, _ := reader.ReadString('\n')
	email = strings.Trim(email, "\r\n")
	fmt.Printf("Sending authentication code to %s\n", email)
	code := authcode.NewAuthCode(email)
	msg := ses_mail.NewAuthMessage(email, code)
	ses_mail.SendMessage(host, email, msg)
	fmt.Printf("Please enter authentication code: ")
	c, _ := reader.ReadString('\n')
	c = strings.Trim(c, "\r\n")
	valid := authcode.ValidateAuthCode(email, c)
	if valid {
		fmt.Println("Valid!")
	} else {
		fmt.Println("Not valid.")
		return
	}
	fmt.Println("Your validation token is below. Please present on your next visit.")
	newToken := authjwt.NewJWT(email, "user")
	newToken.Sign(secret)
	bytes, _ := json.Marshal(newToken)
	fmt.Println(string(bytes))
}
