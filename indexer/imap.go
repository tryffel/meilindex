/*
 * Meilindex - mail indexing and search tool.
 * Copyright (C) 2020 Tero Vierimaa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 *
 */

package indexer

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/jaytaylor/html2text"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"time"
)

type Imap struct {
	Url                 string
	Tls                 bool
	TlsSkipVerification bool
	Username            string
	Password            string

	client  *client.Client
	mailbox *imap.MailboxStatus
}

func (i *Imap) Connect() error {
	var err error
	if i.Tls {
		if i.TlsSkipVerification {
			i.client, err = client.DialTLS(i.Url, &tls.Config{
				InsecureSkipVerify: true,
			})
		} else {
			i.client, err = client.DialTLS(i.Url, nil)
		}
	}

	if err != nil {
		err = fmt.Errorf("connect server: %v", err)
	} else {
		err = i.client.Login(i.Username, i.Password)
		if err != nil {
			err = fmt.Errorf("login: %v", err)
		}
	}

	return err
}

func (i *Imap) Disconnect() error {
	if i.client != nil {
		return i.client.Logout()
	}
	return nil
}

func (i *Imap) Mailboxes() []string {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- i.client.List("", "*", mailboxes)
	}()

	results := make([]string, 0)

	num := 0
	for m := range mailboxes {
		results = append(results, m.Name)
		num += 1
	}

	return results

}

func (i *Imap) SelectMailbox(name string) error {
	mbox, err := i.client.Select(name, true)
	if err != nil {
		return err
	}
	i.mailbox = mbox

	logrus.Infof("Mailbox has %d mails", i.mailbox.Messages)

	return nil
}

func (i *Imap) FetchMail() ([]*Mail, error) {
	messages := make(chan *imap.Message, i.mailbox.Messages)
	done := make(chan error, 1)

	sequence := &imap.SeqSet{}
	start := 1
	stop := i.mailbox.Messages

	//stop = 10

	sequence.AddRange(stop, uint32(start))
	section := &imap.BodySectionName{}

	go func() {
		done <- i.client.Fetch(sequence, []imap.FetchItem{section.FetchItem(), imap.FetchUid}, messages)
	}()
	<-done

	mails := make([]*Mail, len(messages))

	folder := i.mailbox.Name
	if folder == "INBOX" {
		folder = "Inbox"
	}

	for i := 0; i < len(mails); i++ {
		msg := <-messages
		parsed, err := mail.CreateReader(msg.GetBody(section))
		if err != nil {
			logrus.Errorf("parse mail: %v", err)
		}

		m, err := mailToMail(parsed)
		m.Folder = folder

		mails[i] = m
	}

	return mails, nil
}

func mailToMail(m *mail.Reader) (*Mail, error) {
	var err error
	h := m.Header
	date, err := h.Date()
	if err != nil {
		logrus.Errorf("getString date: %v", err)
		date = time.Unix(0, 0)
	} else {

	}
	out := &Mail{
		From:      h.Get("From"),
		To:        h.Get("To"),
		Cc:        h.Get("Cc"),
		Timestamp: date,
		Subject:   h.Get("Subject"),
	}

	out.Uid, err = h.MessageID()

	s, err := h.Subject()
	if err == nil {
		out.Subject = s
	}
	inlineHeaders := 0
	for {
		part, err := m.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			logrus.Errorf("parse mail part: %v", err)
			break
		}

		switch part.Header.(type) {
		case *mail.InlineHeader:
			inlineHeaders += 1
			// accept only 1st inline header
			if inlineHeaders == 1 {
				out.Body, err = html2text.FromReader(part.Body, html2text.Options{
					PrettyTables: false,
				})
			} else {
				b, err := ioutil.ReadAll(part.Body)
				if err != nil {
					logrus.Errorf("read message attachment: %v", err)
				} else {
					out.Attachments = append(out.Attachments, b)
				}
			}

			if err != nil {
				logrus.Errorf("Read html body into text: %v", err)
			}
		case *mail.AttachmentHeader:
			b, err := ioutil.ReadAll(part.Body)
			if err != nil {
				logrus.Errorf("read message attachment: %v", err)
			} else {
				out.Attachments = append(out.Attachments, b)
			}

		}
	}
	return out, err
}
