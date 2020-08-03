/*
Copyright 2020 VMWare

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
 This code is adapted from https://github.com/prometheus/alertmanager/blob/a75cd02786dfecd25e2469fc4df5d920e6b9c226/notify/email/email.go
*/

package smtp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
	"time"

	"github.com/bitnami-labs/kubewatch/config"
	"github.com/mkmik/multierror"
	"github.com/sirupsen/logrus"
)

func sendEmail(conf config.SMTP, msg string) error {
	ctx := context.Background()

	host, port, err := net.SplitHostPort(conf.Smarthost)
	if err != nil {
		return err
	}

	var (
		c       *smtp.Client
		conn    net.Conn
		success = false
	)

	tlsConfig := &tls.Config{}
	if port == "465" {

		if tlsConfig.ServerName == "" {
			tlsConfig.ServerName = host
		}

		conn, err = tls.Dial("tcp", conf.Smarthost, tlsConfig)
		if err != nil {
			return fmt.Errorf("establish TLS connection to server: %w", err)
		}
	} else {
		var (
			d   = net.Dialer{}
			err error
		)
		conn, err = d.DialContext(ctx, "tcp", conf.Smarthost)
		if err != nil {
			return fmt.Errorf("establish connection to server: %w", err)
		}
	}
	c, err = smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer func() {
		// Try to clean up after ourselves but don't log anything if something has failed.
		if err := c.Quit(); success && err != nil {
			logrus.Warnf("failed to close SMTP connection: %v", err)
		}
	}()

	if conf.Hello != "" {
		err = c.Hello(conf.Hello)
		if err != nil {
			return fmt.Errorf("send EHLO command: %w", err)
		}
	}

	// Global Config guarantees RequireTLS is not nil.
	if conf.RequireTLS {
		if ok, _ := c.Extension("STARTTLS"); !ok {
			return fmt.Errorf("'require_tls' is true (default) but %q does not advertise the STARTTLS extension", conf.Smarthost)
		}
		if tlsConfig.ServerName == "" {
			tlsConfig.ServerName = host
		}

		if err := c.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("send STARTTLS command: %w", err)
		}
	}

	if ok, mech := c.Extension("AUTH"); ok {
		auth, err := auth(conf.Auth, host, mech)
		if err != nil {
			return fmt.Errorf("find auth mechanism: %w", err)
		}
		if auth != nil {
			if err := c.Auth(auth); err != nil {
				return fmt.Errorf("%T auth: %w", auth, err)
			}
		}
	}

	addrs, err := mail.ParseAddressList(conf.From)
	if err != nil {
		return fmt.Errorf("parse 'from' addresses: %w", err)
	}
	if len(addrs) != 1 {
		return fmt.Errorf("must be exactly one 'from' address (got: %d)", len(addrs))
	}
	if err = c.Mail(addrs[0].Address); err != nil {
		return fmt.Errorf("send MAIL command: %w", err)
	}
	addrs, err = mail.ParseAddressList(conf.To)
	if err != nil {
		return fmt.Errorf("parse 'to' addresses: %w", err)
	}
	for _, addr := range addrs {
		if err = c.Rcpt(addr.Address); err != nil {
			return fmt.Errorf("send RCPT command: %w", err)
		}
	}

	// Send the email headers and body.
	message, err := c.Data()
	if err != nil {
		return fmt.Errorf("send DATA command: %w", err)
	}
	defer message.Close()

	if conf.Headers == nil {
		conf.Headers = map[string]string{}
	}
	if _, ok := conf.Headers["Subject"]; !ok {
		s := conf.Subject
		if s == "" {
			s = defaultSubject
		}
		conf.Headers["Subject"] = s
	}
	if _, ok := conf.Headers["To"]; !ok {
		conf.Headers["To"] = conf.To
	}
	if _, ok := conf.Headers["From"]; !ok {
		conf.Headers["From"] = conf.From
	}

	buffer := &bytes.Buffer{}
	for header, value := range conf.Headers {
		fmt.Fprintf(buffer, "%s: %s\r\n", header, mime.QEncoding.Encode("utf-8", value))
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	if _, ok := conf.Headers["Message-Id"]; !ok {
		fmt.Fprintf(buffer, "Message-Id: %s\r\n", fmt.Sprintf("<%d.%d@%s>", time.Now().UnixNano(), rand.Uint64(), hostname))
	}

	multipartBuffer := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(multipartBuffer)

	fmt.Fprintf(buffer, "Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	fmt.Fprintf(buffer, "Content-Type: multipart/alternative;  boundary=%s\r\n", multipartWriter.Boundary())
	fmt.Fprintf(buffer, "MIME-Version: 1.0\r\n\r\n")

	_, err = message.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("write headers: %w", err)
	}
	w, err := multipartWriter.CreatePart(textproto.MIMEHeader{
		"Content-Transfer-Encoding": {"quoted-printable"},
		"Content-Type":              {"text/plain; charset=UTF-8"},
	})
	if err != nil {
		return fmt.Errorf("create part for text template: %w", err)
	}

	qw := quotedprintable.NewWriter(w)
	_, err = qw.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("write text part: %w", err)
	}
	err = qw.Close()
	if err != nil {
		return fmt.Errorf("close text part: %w", err)
	}

	err = multipartWriter.Close()
	if err != nil {
		return fmt.Errorf("close multipartWriter: %w", err)

	}

	_, err = message.Write(multipartBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("write body buffer: %w", err)
	}

	log.Printf("sending via %s:%s, to: %q, from: %q : %s ", host, port, conf.To, conf.From, msg)
	return nil
}

func auth(conf config.SMTPAuth, host, mechs string) (smtp.Auth, error) {
	username := conf.Username

	// If no username is set, keep going without authentication.
	if username == "" {
		logrus.Debugf("smtp_auth_username is not configured. Attempting to send email without authenticating")
		return nil, nil
	}

	var errs []error
	for _, mech := range strings.Split(mechs, " ") {
		switch mech {
		case "CRAM-MD5":
			secret := string(conf.Secret)
			if secret == "" {
				errs = append(errs, fmt.Errorf("missing secret for CRAM-MD5 auth mechanism"))
				continue
			}
			return smtp.CRAMMD5Auth(username, secret), nil

		case "PLAIN":
			password := string(conf.Password)
			if password == "" {
				errs = append(errs, fmt.Errorf("missing password for PLAIN auth mechanism"))
				continue
			}
			identity := conf.Identity

			return smtp.PlainAuth(identity, username, password, host), nil
		case "LOGIN":
			password := string(conf.Password)
			if password == "" {
				errs = append(errs, fmt.Errorf("missing password for LOGIN auth mechanism"))
				continue
			}
			return LoginAuth(username, password), nil
		}
	}
	if len(errs) == 0 {
		errs = append(errs, fmt.Errorf("unknown auth mechanism: %q", mechs))
	}
	return nil, multierror.Join(errs)
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

// Used for AUTH LOGIN. (Maybe password should be encrypted)
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch strings.ToLower(string(fromServer)) {
		case "username:":
			return []byte(a.username), nil
		case "password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unexpected server challenge")
		}
	}
	return nil, nil
}
