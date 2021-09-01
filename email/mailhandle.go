
package kpnmmail

import (
	tls   "crypto/tls"

	gomail "gopkg.in/gomail.v2"
)


type Email struct{
	host string
	port int
	address string
	password string
	dialer *gomail.Dialer
	sender gomail.SendCloser
}

func NewEmail(host string, port int, address string, password string)(mail *Email){
	mail = &Email{
		host: host,
		port: port,
		address: address,
		password: password,
		dialer: gomail.NewDialer(host, port, address, password),
	}
	mail.dialer.TLSConfig = &tls.Config{ InsecureSkipVerify: true }
	return
}

func (mail *Email)Login()(err error){
	mail.sender, err = mail.dialer.Dial()
	return
}

func (mail *Email)Close()(error){
	return mail.sender.Close()
}

func (mail *Email)SendMail(to string, title string, text string)(err error){
	msg := gomail.NewMessage()
	msg.SetHeader("From", mail.address)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", title)
	msg.SetBody("text/html", text)
	err = mail.sender.Send(mail.address, []string{to}, msg)
	if err != nil {
		return
	}
	return nil
}

func (mail *Email)SendHtml(to string, title string, path string, value interface{})(err error){
	text, err := ExeHtmlTemp(path, value)
	if err != nil {
		return err
	}
	return mail.SendMail(to, title, text)
}


