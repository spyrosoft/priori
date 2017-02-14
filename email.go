package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"net/smtp"
)

func sendEmail(to string, subject string, message string) (err error) {
	err = sendEmailFrom(to, siteData.NoReplyAddressName, siteData.NoReplyAddress, subject, message)
	return
}

func sendEmailFrom(to string, fromName string, from string, subject string, message string) (err error) {
	if !siteDataLoaded {
		err = errors.New("Outgoing email credentials have not been set. Cannot send message.")
		return
	}

	fromHeader := mail.Address{fromName, from}

	headers := make(map[string]string)
	headers["From"] = fromHeader.String()
	headers["To"] = to
	headers["Subject"] = subject

	rawMessage := ""
	for headerName, headerValue := range headers {
		rawMessage += fmt.Sprintf("%s: %s\r\n", headerName, headerValue)
	}
	rawMessage += "\r\n" + message

	mailAuth := smtp.PlainAuth("", siteData.NoReplyAddress, siteData.NoReplyPassword, siteData.Host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         siteData.Host,
	}

	tcpConnection, err := tls.Dial("tcp", siteData.Host+":"+siteData.Port, tlsConfig)
	if err != nil {
		return
	}

	smtpClient, err := smtp.NewClient(tcpConnection, siteData.Host)
	if err != nil {
		return
	}

	defer smtpClient.Quit()

	err = smtpClient.Auth(mailAuth)
	if err != nil {
		return
	}

	err = smtpClient.Mail(siteData.NoReplyAddress)
	if err != nil {
		return
	}

	err = smtpClient.Rcpt(to)
	if err != nil {
		return
	}

	emailStream, err := smtpClient.Data()
	if err != nil {
		return
	}

	_, err = emailStream.Write([]byte(rawMessage))
	if err != nil {
		return
	}

	err = emailStream.Close()
	if err != nil {
		return
	}

	return
}

func apiNotifyAdmin(r *http.Request) (results string) {
	message := r.PostFormValue("message")
	if r.PostFormValue("token") != "ytMOJPatwt" || message == "" {
		return
	}
	err := notifyAdmin(message)
	if err != nil {
		resultsBytes, _ := json.Marshal(ErrorJSON{Errors: []string{err.Error()}})
		results = string(resultsBytes)
	} else {
		resultsBytes, _ := json.Marshal(SuccessJSON{Success: true})
		results = string(resultsBytes)
	}
	return
}

func notifyAdmin(message string) (err error) {
	err = sendEmail(siteData.AdminEmail, "An Unexpected Error Occurred - Priori", message)
	return
}
