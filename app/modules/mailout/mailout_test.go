package mailout

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const EXPECTED_MAIL_BODY = `Subject: Buchungskalender (Belegte Tage: 100)


Der aktuelle Buchungskalender zur Ferienwohnung Strandsommer E10.
Belegte Tage: 100

*** Belegungsplan ***
`

func TestMailOut(t *testing.T) {
	body := MailBody{Days: 100}

	buf, err := createMail(body)

	if err != nil {
		t.Errorf("Error: %s", err)
	}
	if buf.String() == "" {
		t.Errorf("Error: %s", "Mail body is empty")
	}
	assert.Equal(t, EXPECTED_MAIL_BODY, buf.String())
}
