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
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Mail struct {
	// Uid is hash calculated from id. Uid contains only ascii characters.
	Uid string `json:"uid"`
	// Original message id
	Id              string    `json:"id"`
	From            string    `json:"from"`
	To              []string  `json:"to"`
	Cc              []string  `json:"cc"`
	Subject         string    `json:"subject"`
	Body            string    `json:"body"`
	Timestamp       time.Time `json:"date"`
	Folder          string    `json:"folder"`
	Attachments     [][]byte
	AttachmentNames []string
}

func (m *Mail) String() string {
	return fmt.Sprintf(
		`
id: %s,
folder: %s
date: %s
from: %s,
to: %s, 
cc: %s,
subject: %s,

%s

`, m.Id, m.Folder, m.DateTime(), m.From, m.To, m.Cc, m.Subject, m.Body)
}

// Date returns date part of timestamp
func (m *Mail) Date() string {
	return m.Timestamp.Format("2006-01-02")
}

// DateTime returns date and local time
func (m *Mail) DateTime() string {
	return m.Timestamp.Format("2006-01-02 15:04")
}

// ShortDateTime returns short / simple format for date, today, yesterday at xx.xx, date
func (m *Mail) ShortDateTime() string {
	today := time.Now()
	distance := today.Sub(m.Timestamp)
	minutes := distance.Minutes()
	hours := minutes / 60

	if today.Day() == m.Timestamp.Day() && hours < 24 {
		return fmt.Sprintf("Today, %d:%d", m.Timestamp.Hour(), m.Timestamp.Minute())
	}

	if today.Day() == m.Timestamp.Day()+1 && hours < 48 {
		return fmt.Sprintf("Yesterday, %d:%d", m.Timestamp.Hour(), m.Timestamp.Minute())
	}

	return m.Date()
}

// Sanitize makes various mail attributes nicer to read.
func (m *Mail) Sanitize() {

}

func (m *Mail) HighlightedBody() string {
	body := m.Body
	if strings.Contains(body, "<em>") {
		body = strings.Replace(body, "<em>", "[black:orange:]", -1)
		body = strings.Replace(body, "</em>", "[-:-:]", -1)
	}
	return body
}

func (m *Mail) HighlightedSubject() string {
	subject := m.Subject
	if strings.Contains(subject, "<em>") {
		subject = strings.Replace(subject, "<em>", "[black:orange:]", -1)
		subject = strings.Replace(subject, "</em>", "[-:-:]", -1)
	}
	return subject
}

func (m *Mail) HighlightedFrom() string {
	from := m.From
	if strings.Contains(from, "<em>") {
		from = strings.Replace(from, "<em>", "[black:orange:]", -1)
		from = strings.Replace(from, "</em>", "[-:-:]", -1)
	}
	return from
}

var addressNames = regexp.MustCompile(`\"([^'\"]+)\"\s<([\w.]+@[a-zA-Z.]+)>`)
var plainAddress = regexp.MustCompile(`([\w.]+@[a-zA-Z.]+)`)
var escapedNames = regexp.MustCompile(`\"'([^\"]+)'\"\s<([\w.]+@[a-zA-Z.]+)>`)

// strip multiple addresses and possible names to names-only list
func stripdAddressNames(address string) []string {
	var out []string
	// Catch "person" <address>
	matches := addressNames.FindAllStringSubmatch(address, -1)
	if len(matches) > 0 {
		for _, v := range matches {
			out = append(out, v[1])
		}
		return out
	}
	// catch "'sender'" <address>
	matches = escapedNames.FindAllStringSubmatch(address, -1)
	if len(matches) > 0 {
		for _, v := range matches {
			out = append(out, v[1])
		}
		return out
	}

	// catch address
	matches = plainAddress.FindAllStringSubmatch(address, -1)
	if len(matches) > 0 {
		for _, v := range matches {
			out = append(out, v[1])
		}
		return out
	}
	return []string{address}
}
