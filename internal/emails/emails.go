package emails

import (
	"sync"
	"tidal.lol/internal/smtpd"
	"time"
)

var (
	// Emails is a list of all emails sent to the server
	Emails []smtpd.Email

	mu = sync.Mutex{}
)

// Get returns all emails sent to a given address
func Get(address string) []smtpd.Email {
	mu.Lock()
	defer mu.Unlock()

	var emails []smtpd.Email
	for _, email := range Emails {
		if email.To == address {
			emails = append(emails, email)
		}
	}

	return emails
}

// Insert inserts an email into the database
func Insert(email smtpd.Email) {
	mu.Lock()
	defer mu.Unlock()

	Emails = append(Emails, email)
}

// Delete deletes all emails sent to a given address
func Delete(uniqueID string) {
	mu.Lock()
	defer mu.Unlock()

	for index, value := range Emails {
		if value.UniqueID == uniqueID {
			Emails = append(Emails[:index], Emails[index+1:]...)
		}
	}
}

// DeleteAfter deletes all emails sent to a given address after a given duration
func DeleteAfter(uniqueID string, duration time.Duration) {
	time.AfterFunc(duration, func() {
		mu.Lock()
		defer mu.Unlock()

		for index, value := range Emails {
			if value.UniqueID == uniqueID {
				Emails = Emails[index:]
			}
		}
	})
}
