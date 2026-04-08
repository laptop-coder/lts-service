package service

import (
	"backend/pkg/env"
	"backend/pkg/logger"
	"bytes"
	"context"
	"crypto/tls"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
)

//go:embed email_templates/*.html
var templatesFS embed.FS

type EmailService interface {
	Send(ctx context.Context, to []string, subject, body string) error
	SendInviteLink(ctx context.Context, to *string, link string) error
}

type emailService struct {
	config    EmailServiceConfig
	log       logger.Logger
	templates *template.Template
}

func NewEmailService(
	config EmailServiceConfig,
	log logger.Logger,
) (EmailService, error) {
	tmpl, err := template.ParseFS(templatesFS, "email_templates/*.html")
	if err != nil {
		return nil, err
	}
	return &emailService{
		config:    config,
		log:       log,
		templates: tmpl,
	}, nil
}

func (s *emailService) Send(ctx context.Context, to []string, subject, body string) error {
	// SMTP server address
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	// Usual connection
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()
	// Assemble TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: env.GetBoolRequired("SMTP_SKIP_VERIFY"), // TODO: move to the config
		ServerName:         s.config.Host,
	}
	// STARTTLS: upgrade connection to TLS
	if err := conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}
	// Auth
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("failed to auth: %w", err)
	}
	// Sender
	if err := conn.Mail(s.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	// Recipients
	for _, recipient := range to {
		if err := conn.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}
	}
	// Get data writer
	w, err := conn.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer w.Close()
	// Assemble email
	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		s.config.From, strings.Join(to, ", "), subject, body,
	))
	// Send email
	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	return nil
}

func (s *emailService) SendInviteLink(ctx context.Context, to *string, link string) error {
	if to == nil {
		return fmt.Errorf("email cannot be nil")
	} else if strings.TrimSpace(*to) == "" {
		return fmt.Errorf("email cannot be empty or only whitespace")
	}
	var body bytes.Buffer
	if err := s.templates.ExecuteTemplate(&body, "invite_link.html", map[string]string{
		"InviteLink": link,
		"FooterText": env.GetStringRequired("FOOTER_TEXT"), // TODO: move to the config
	}); err != nil {
		return err
	}
	return s.Send(ctx, []string{*to}, "Добро пожаловать!", body.String())
}
