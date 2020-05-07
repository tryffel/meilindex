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
	"flag"
	"github.com/sirupsen/logrus"
)

var meilisearchHost = "http://localhost:7700"
var meilisearchIndex = "mail"
var meilisearchApiKey = "masterKey"

func main() {
	index := flag.Bool("index", false, "Index mail")
	query := flag.String("query", "", "Query to search")

	flag.Parse()

	if *index {
		indexMail()
	} else if *query != "" {
		searchMail(*query)
	}
}

func indexMail() {

	client := &Imap{
		Url:                 "imap.mymail.com:993",
		Tls:                 true,
		TlsSkipVerification: false,
		Username:            "me@mymail.com",
		Password:            "memailing",
	}

	err := client.Connect()
	if err != nil {
		logrus.Error(err)
		return
	} else {
		logrus.Info("Logged in")
	}

	defer client.Disconnect()

	err = client.SelectMailbox("INBOX")
	if err != nil {
		logrus.Errorf("select mailbox: %v", err)
		return
	}

	mails, err := client.FetchMail()
	ms := Meilisearch{
		Url:    meilisearchHost,
		Index:  meilisearchIndex,
		ApiKey: meilisearchApiKey,
	}

	err = ms.Connect()
	if err != nil {
		logrus.Error(err)
	}

	err = ms.IndexMail(mails)
	if err != nil {
		logrus.Error(err)
	}

}

func searchMail(query string) {
	ms := Meilisearch{
		Url:    meilisearchHost,
		Index:  meilisearchIndex,
		ApiKey: meilisearchApiKey,
	}

	err := ms.Connect()
	if err != nil {
		logrus.Error(err)
	}

	err = ms.Query(query)
}
