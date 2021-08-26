
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
}

func NewEmail(host string, port int, address string, password string)(mail *Email){
	mail = new(Email)
	mail.host = host
	mail.port = port
	mail.address = address
	mail.password = password
	return mail
}

func (mail *Email)SendMail(to string, title string, text string)(err error){
	msg := gomail.NewMessage()
	msg.SetHeader("From", mail.address)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", title)
	msg.SetBody("text/html", text)
	dial := gomail.NewDialer(mail.host, mail.port, mail.address, mail.password)
	dial.TLSConfig = &tls.Config{ InsecureSkipVerify: true }
	err = dial.DialAndSend(msg)
	if err != nil {
		return err
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


