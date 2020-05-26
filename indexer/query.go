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
	"github.com/meilisearch/meilisearch-go"
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"tryffel.net/go/meilindex/config"
)

var queryPattern = `([+-])?([a-zA-Z]+):(\w+|'[\w ]+')`
var queryRegex = regexp.MustCompile(queryPattern)

func SearchMail(query string, filter string) {
	ms := Meilisearch{
		Url:    config.Conf.Meilisearch.Url,
		Index:  config.Conf.Meilisearch.Index,
		ApiKey: config.Conf.Meilisearch.ApiKey,
	}

	err := ms.Connect()
	if err != nil {
		logrus.Error(err)
	}

	mails, _, err := ms.Query(query, filter)
	if err != nil {
		logrus.Error(err)
	} else {
		yellow := ansi.ColorCode("yellow+i:black")
		reset := ansi.ColorCode("reset")
		for i := 0; i < len(mails); i++ {
			mail := mails[i]
			if strings.Contains(mail.Body, "<em>") {
				mail.Body = strings.Replace(mail.Body, "<em>", yellow, -1)
				mail.Body = strings.Replace(mail.Body, "</em>", reset, -1)
			}

			mail.Subject = ansi.Blue + mail.Subject + ansi.Reset
			fmt.Println("-----------------")
			fmt.Printf("%d %s", i, mail.String())
		}
	}
}

func (m *Meilisearch) Query(query, filter string) ([]*Mail, int, error) {

	//yellow := ansi.ColorCode("yellow+i:black")
	//reset := ansi.ColorCode("reset")

	res, err := m.client.Search(m.Index).Search(meilisearch.SearchRequest{
		Query:                 query,
		Limit:                 40,
		AttributesToCrop:      []string{"message:200"},
		AttributesToHighlight: []string{"message"},
		Filters:               filter,
	})

	if err != nil {
		return nil, -1, err
	}

	result := make([]*Mail, len(res.Hits))

	for i, v := range res.Hits {
		isMap, ok := v.(map[string]interface{})
		formatted := isMap["_formatted"]
		mail := &Mail{}
		if ok {
			if isFormatted, ok := formatted.(map[string]interface{}); ok {
				mail.Body = get("message", isFormatted)
			}
		} else {
			mail.Body = get("message", isMap)
		}
		mail.Id = get("uid", isMap)
		mail.From = get("from", isMap)
		mail.To = get("to", isMap)
		mail.Cc = get("cc", isMap)
		mail.Subject = get("subject", isMap)
		mail.Folder = get("folder", isMap)
		mail.Date = get("date", isMap)

		result[i] = mail

	}
	//println(mail.String())
	//println("\n\n")
	//p//rintln("===============================================")
	return result, int(res.ProcessingTimeMs), nil
}

func get(key string, container map[string]interface{}) string {
	val, ok := container[key].(string)
	if !ok {
		return ""
	}
	return val
}
