package emails

import (
	"net"
	"tidal.lol/internal/database"
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
				Str("to", to[0]).
				Str("unique_id", email.UniqueID).
				Msg("Received email")
		}

		query, err := database.DB.Query("SELECT * FROM mailbox WHERE email=$email", map[string]any{"email": to[0]})
		if err != nil {
			logging.Logger.Error().
				Err(err).
				Str("email", to[0]).
				Msg("Failed to query mailbox")
		}

		// If the mailbox doesn't exist, send it to the temporary inbox
		if len(query.([]interface{})[0].(map[string]interface{})["result"].([]interface{})) == 0 {
			Insert(*email)
			go DeleteAfter(email.UniqueID, time.Hour)
		} else { // Otherwise, send it to the inbox
			_, err := database.DB.Create("inbox", map[string]any{
				"token": query.([]interface{})[0].(map[string]interface{})["result"].([]interface{})[0].(map[string]interface{})["token"],

				"to":   to[0],
				"from": from,

				"subject": email.Subject,
				"date":    email.Date,

				"body": map[string]any{
					"text": email.Body.Text,
					"html": email.Body.HTML,
				},
			})
			if err != nil {
				logging.Logger.Error().
					Err(err).
					Str("email", to[0]).
					Msg("Failed to insert into inbox")
			}
		}

		return nil
	}
}
