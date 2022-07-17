package main

import (
	"auth/pkg/authmail"
	"auth/pkg/authserver"
	"os"
	"time"
)

var secret []byte = []byte("test secret")

func main() {
	srv := &authserver.Server{
		Address: "localhost",
		Port:    8080,
		SESHost: authmail.Host{
			Username: os.Getenv("SES_USERNAME"),
			Password: os.Getenv("SES_PASSWORD"),
			Host:     "email-smtp.us-east-2.amazonaws.com",
			Port:     587,
			Sender:   "jnichols2719@protonmail.com",
		},
	}
	srv.Start()
	time.Sleep(time.Minute * 2)
	srv.Stop()
	time.Sleep(time.Second * 2)
	/*
		rand.Seed(time.Now().UnixNano())
		host := authmail.Host{
			Username: os.Getenv("SES_USERNAME"),
			Password: os.Getenv("SES_PASSWORD"),
			Host:     "email-smtp.us-east-2.amazonaws.com",
			Port:     587,
			Sender:   "jnichols2719@protonmail.com",
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Please enter authentication token: ")
		token, _ := reader.ReadString('\n')
		if token = strings.Trim(token, "\r\n"); token != "" {
			jwt, valid, _ := authjwt.Verify(token, secret)
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
		msg := authmail.NewAuthMessage(email, code.Code)
		authmail.SendMessage(host, email, msg)
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
		newToken := authjwt.NewJWT(email, map[string]interface{}{"user-type": "user"})
		output, err := authjwt.Export(newToken, secret)
		if err != nil {
			panic(err)
		}
		fmt.Println(output)
	*/
}
