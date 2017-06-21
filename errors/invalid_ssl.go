package errors

import "fmt"

type InvalidSSLCert struct {
	URL    string
	Reason string
}

func NewInvalidSSLCert(url, reason string) *InvalidSSLCert {
	return &InvalidSSLCert{
		URL:    url,
		Reason: reason,
	}
}

func (err *InvalidSSLCert) Error() string {
	message := fmt.Sprintf("Received invalid SSL certificate from %s", err.URL)
	if err.Reason != "" {
		message += " - " + err.Reason
	}
	return message
}
