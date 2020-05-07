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

package main

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/mail"
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
		parsed, err := mail.ReadMessage(msg.GetBody(section))
		if err != nil {
			logrus.Errorf("parse mail: %v", err)
			continue
		}

		h := parsed.Header
		m := &Mail{
			Id:      int64(msg.Uid),
			From:    h.Get("From"),
			To:      h.Get("To"),
			Cc:      h.Get("Cc"),
			Date:    h.Get("Date"),
			Subject: h.Get("Subject"),
			Folder:  folder,
		}

		body, err := ioutil.ReadAll(parsed.Body)
		m.Body = string(body)

		mails[i] = m
	}

	return mails, nil
}
