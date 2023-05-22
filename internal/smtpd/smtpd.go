package smtpd

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
	"github.com/mhale/smtpd"
	"net"
	"net/mail"
	"time"
)

// NewServer creates a new SMTP listener
func NewServer(port int, handler func(origin net.Addr, from string, to []string, data []byte) error) error {
	server := &smtpd.Server{
		Appname: "tidal",
		Addr:    fmt.Sprintf("0.0.0.0:%v", port),

		Handler: handler,
		HandlerRcpt: func(remoteAddr net.Addr, from string, to string) bool {
			return true
		},
	}

	return server.ListenAndServe()
}

// Parse parses an email from a byte array
func Parse(from string, to []string, d []byte) (*Email, error) {
	email, err := mail.ReadMessage(bytes.NewReader(d))
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(d)
	env, err := enmime.ReadEnvelope(reader)
	if err != nil {
		return nil, err
	}

	return &Email{
		UniqueID: uuid.NewString(),

		To:   to[0],
		From: from,

		Date:    time.Now().Format(time.RFC3339),
		Subject: email.Header.Get("Subject"),

		Body: struct {
			Text string `json:"text"`
			HTML string `json:"html"`
		}(struct {
			Text string
			HTML string
		}{Text: env.Text, HTML: env.HTML}),
	}, nil
}
