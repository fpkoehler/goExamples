/* from https://gist.github.com/chrisgillis/10888032 */

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

// SSL/TLS Email Example

func main() {

//	from := mail.Address{"", "fkoehler@aol.com"}
//	to := mail.Address{"", "fred@fpkoehler.com"}
	from := mail.Address{"", "fred@fpkoehler.com"}
	to := mail.Address{"", "fkoehler@aol.com"}
	subj := "Go email"
	body := "This is an example body.\n With two lines."

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
//	headers["To"] = to.String()
	headers["To"] = "fred@fpkoehler.com, fkoehler@aol.com"
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
//	servername := "smtp.aol.com:465"
	servername := "smtp.1and1.com:465"

	host, _, _ := net.SplitHostPort(servername)

//	auth := smtp.PlainAuth("", "fkoehler@aol.com", "taxi11", host)
	auth := smtp.PlainAuth("", "fred@fpkoehler.com", "Analog==", host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Panic(err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// SECOND RECEPIENT!!
	if err = c.Rcpt("fred@fpkoehler.com"); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

}
