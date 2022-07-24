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

// Generate an smtp.PlainAuth based on a Host object.
//
// Input:
//   - host Host: Filled host data
// Output:
//   - smtp.Auth: Used in smtp.SendMail to authenticate with the SMTP provider.
func (host Host) plainAuth() smtp.Auth {
	return smtp.PlainAuth("", host.Username, host.Password, host.Host)
}

// Generate a string address for contacting the SMTP host.
//
// Input:
//   - host Host: Filled host data
// Output:
//   - string: SMTP host address
func (host Host) address() string {
	return fmt.Sprintf("%s:%d", host.Host, host.Port)
}

// Generate a new authentication message given a target email and a gate code.
//
// Input:
//   - sendTo string: Email to send the code to
//   - code string: Authentication code
// Output:
//   - []byte: Properly formatted message for sending through smtp.
func NewAuthMessage(sendTo string, code string) []byte {
	msg := fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: Authentication Code\r\n"+
			"\r\n"+
			"Your authentication code is %s.\n"+
			"This code will expire in 5 minutes.\r\n",
		sendTo, code,
	)
	return []byte(msg)
}

// Send a message from a Host to an email address.
//
// Input:
//   - sendFrom Host: SMTP host information
//   - sendTo string: Target email address
//   - msg []byte: Message to send. Use NewAuthMessage() to generate this.
// Output:
//   - error: Any error that occurs when sending mail. This only includes failure to *access* the SMTP server,
//   or invalid credentials; if all configuration is correct but the email fails to go through or sendTo is not
//   a valid address, no error will be returned.
func SendMessage(sendFrom Host, sendTo string, msg []byte) error {
	auth := sendFrom.plainAuth()
	target := []string{sendTo}
	addr := sendFrom.address()
	return smtp.SendMail(addr, auth, sendFrom.Sender, target, msg)
}
