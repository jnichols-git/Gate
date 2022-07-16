# pkg/authmail

The `authmail` package gives email functionality to `auth`.

## Usage

`authmail` exposes:

- the `Host` struct: defines smtp hosts
- `NewAuthMessage(to, code string)`: generates the highly-formatted `[]byte` used in `net/smtp`
- `SendMessage(from Host, to string, msg []byte)`: sends a message using smtp

Generally, the control flow will look something like this (in conjuction with `authcode`) to verify a user.

```
// Create SMTP host
host := authmail.Host{
	Username: os.Getenv("SES_USERNAME"),
	Password: os.Getenv("SES_PASSWORD"),
	Host:     "email-smtp.us-east-2.amazonaws.com",
	Port:     587,
	Sender:   "example@domain.com",
}

// Generate a code for this user
sendToEmail := "someuser@gmail.com"
code := authcode.NewAuthCode(sendToEmail)
// Send the code as an email
msg := authmail.NewAuthMessage(sendToEmail, code.Code)
authmail.SendMessage(host, sendToEmail, msg)

// In this example, assume we wait until the user inputs the below...
inputCode := <user input>
// Validate the authentication code
valid := authcode.ValidateAuthCode(inputCode)
if valid {
    fmt.Println("Hooray! Validated!")
} else {
    fmt.Println("Invalid code :(")
}
```

## Why SMTP?

SMTP is, from what I have seen, widely used, relatively easy to set up, and highly configurable.
