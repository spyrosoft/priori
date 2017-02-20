package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"net/smtp"
	runtimeDebug "runtime/debug"
)

func sendEmail(to string, subject string, message string) (err error) {
	err = sendEmailFrom(to, siteData.NoReplyAddressName, siteData.NoReplyAddress, subject, message)
	return
}

func sendEmailFrom(to string, fromName string, from string, subject string, message string) (err error) {
	err = errors.New("Outgoing email credentials have not been set. Cannot send message.")

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

func apiNotifyAdmin(w http.ResponseWriter, r *http.Request) interface{} {
	message := r.PostFormValue("message")
	if r.PostFormValue("token") != "PpPub4GjM4" || message == "" {
		return apiResponse{}
	}
	err := notifyAdmin(message)
	if err != nil {
		return apiResponse{Errors: []string{err.Error()}}
	}
	return apiResponse{Success: true}
}

func notifyAdminResponse(message string, err error) (response apiResponse) {
	err = errors.New(err.Error() + "\n" + string(runtimeDebug.Stack()))
	response = apiResponse{
		Errors: []string{message, adminNotifiedMessage},
		Debug:  []string{err.Error()},
	}
	responseBytes, _ := json.Marshal(response)
	notifyAdmin(string(responseBytes))
	return
}

func notifyAdmin(message string) (err error) {
	err = sendEmail(siteData.AdminEmail, "An Unexpected Error Occurred - Priori", message)
	return
}
