package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"net/smtp"
	"github.com/skiesel/designatedbrewer/utils"
)

type messageRequest struct {
	Subject string
	Message string
}


func AlertMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var m messageRequest
	err := decoder.Decode(&m)
	if err != nil {
		panic(err)
	}

	from := mail.Address{"", config.GetSetting("smtp user")}
	to := mail.Address{"", config.GetSetting("contact address")}
	subj := m.Subject
	body := m.Message

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	host := config.GetSetting("smtp server")
	port := config.GetSetting("smtp port")
	fullHost := host + ":" + port


	auth := smtp.PlainAuth("", config.GetSetting("smtp user"), config.GetSetting("smtp password"), host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", fullHost, tlsconfig)
	if err != nil {
		panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		panic(err)
	}

	// Data
	wr, err := c.Data()
	if err != nil {
		panic(err)
	}

	_, err = wr.Write([]byte(message))
	if err != nil {
		panic(err)
	}

	err = wr.Close()
	if err != nil {
		panic(err)
	}

	c.Quit()

	fmt.Fprint(w, "success")
}
