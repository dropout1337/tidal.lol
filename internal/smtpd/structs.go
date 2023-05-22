package smtpd

// Email represents an email
type Email struct {
	UniqueID string `json:"unique_id"`

	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`

	Date string `json:"date"`

	Body struct {
		Text string `json:"text"`
		HTML string `json:"html"`
	} `json:"body"`
}
