package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"

	"github.com/BurntSushi/toml"
)

type EmailReleaseHandlerConfig struct {
	SmtpServer   string `toml:"smtp_server"`
	SmtpPort     int    `toml:"smtp_port"`
	SmtpUsername string `toml:"smtp_username"`
	SmtpPassword string `toml:"smtp_password"`

	From     string `toml:"from"`
	To       string `toml:"to"`
	Template string `toml:"template"`
}

type EmailReleaseHandler struct {
	smtpServer string
	smtpPort   int
	smtpAuth   smtp.Auth
	from       *mail.Address
	to         *mail.Address
	template   *template.Template
}

func (handler EmailReleaseHandler) Handle(repo *Repo, notification GithubNotification) error {
	var err error

	contents := new(bytes.Buffer)

	contents.Write([]byte(fmt.Sprintf("To: %s\r\n", handler.to)))
	contents.Write([]byte(fmt.Sprintf("From: %s\r\n", handler.from)))
	contents.Write([]byte(fmt.Sprintf("Subject: %s\r\n", notification.Release.Name)))
	contents.Write([]byte("Content-Type: text/html; charset=UTF-8\r\n"))
	contents.Write([]byte("\r\n"))

	// TODO: Pass the Repo into the template too? Or set some stuff like app
	//       name that's in Repo manually into notification?
	if err = handler.template.Execute(contents, notification); err != nil {
		return err
	}

	if err = smtp.SendMail(fmt.Sprintf("%s:%d", handler.smtpServer, handler.smtpPort),
		handler.smtpAuth, handler.from.Address, []string{handler.to.Address},
		contents.Bytes()); err != nil {
		return err
	}

	return nil
}

func NewEmailReleaseHandler(configPrimitive toml.Primitive) (NotificationHandler, error) {
	var err error

	var config EmailReleaseHandlerConfig
	if err = toml.PrimitiveDecode(configPrimitive, &config); err != nil {
		return nil, err
	}

	auth := smtp.PlainAuth("", config.SmtpUsername, config.SmtpPassword, config.SmtpServer)

	var template *template.Template
	if template, err = template.ParseFiles(config.Template); err != nil {
		return nil, err
	}

	var (
		to   *mail.Address
		from *mail.Address
	)

	if to, err = mail.ParseAddress(config.To); err != nil {
		return nil, err
	}

	if from, err = mail.ParseAddress(config.From); err != nil {
		return nil, err
	}

	return &EmailReleaseHandler{config.SmtpServer, config.SmtpPort, auth, from, to, template}, nil
}
