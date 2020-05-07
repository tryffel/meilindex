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
	"fmt"
	"github.com/jaytaylor/html2text"
	"github.com/sirupsen/logrus"
)

type Mail struct {
	Id      int64  `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Cc      string `json:"cc"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Date    string `json:"date"`
	Folder  string `json:"folder"`
}

func (m *Mail) String() string {
	return fmt.Sprintf(
		`
id: %d,
folder: %s
date: %s
from: %s,
to: %s, 
cc: %s,
subject: %s,

%s
`, m.Id, m.Folder, m.Date, m.From, m.To, m.Cc, m.Subject, m.Body)
}

func (m *Mail) BodyPlainText() string {
	text, err := html2text.FromString(m.Body, html2text.Options{
		PrettyTables: false,
	})

	if err != nil {
		logrus.Errorf("get plain text from html body: %v", err)
		return m.Body
	}

	return text
}
