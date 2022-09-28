package emails

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"
)

const tmpl = "MIME-Version: 1.0\r\n" +
	"To: %s\r\n" +
	"Subject: %s\r\n" +
	"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
	"\r\n" +
	"%s\r\n"

type Sender interface {
	Send(string, string, string) error
}

type User struct {
	Login    string
	Password string
}

type sender struct {
	from string
	host string
	auth smtp.Auth
}

func NewSender(host string, port int, user *User, from string) Sender {
	auth := smtp.PlainAuth("", user.Login, user.Password, host)

	return &sender{
		auth: auth,
		host: host + ":" + strconv.Itoa(port),
		from: from,
	}
}

func (s *sender) Send(to string, sub string, msg string) error {
	mail := fmt.Sprintf(tmpl, to, sub, msg)

	c, err := smtp.Dial(s.host)
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{
			ServerName:         s.host,
			InsecureSkipVerify: true,
		}

		if err = c.StartTLS(config); err != nil {
			return err
		}
	}

	if err = c.Auth(s.auth); err != nil {
		return err
	}

	if err = c.Mail(s.from); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(mail))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
