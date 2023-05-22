package emails

import (
	"net"
	"strings"
	"tidal.lol/internal/logging"
	"tidal.lol/internal/smtpd"
	"time"
)

// CreateHandler creates a new handler
func CreateHandler(debug bool) func(origin net.Addr, from string, to []string, data []byte) error {
	// Handler handles incoming emails
	return func(origin net.Addr, from string, to []string, data []byte) error {
		email, err := smtpd.Parse(from, to, data)
		if err != nil {
			return err
		}

		if debug {
			logging.Logger.Info().
				Str("from", from).
				Str("to", strings.Join(to, ", ")).
				Str("unique_id", email.UniqueID).
				Msg("Received email")
		}

		Insert(*email)
		go DeleteAfter(email.UniqueID, time.Hour)

		return nil
	}
}
