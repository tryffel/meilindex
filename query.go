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
	"github.com/meilisearch/meilisearch-go"
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var queryPattern = `([+-])?([a-zA-Z]+):(\w+|'[\w ]+')`
var queryRegex = regexp.MustCompile(queryPattern)

func searchMail(query string, filter string) {
	ms := Meilisearch{
		Url:    meilisearchHost,
		Index:  meilisearchIndex,
		ApiKey: meilisearchApiKey,
	}

	err := ms.Connect()
	if err != nil {
		logrus.Error(err)
	}

	err = ms.Query(query, filter)
}

func (m *Meilisearch) Query(query, filter string) error {

	yellow := ansi.ColorCode("yellow+i:black")
	reset := ansi.ColorCode("reset")

	res, err := m.client.Search(m.Index).Search(meilisearch.SearchRequest{

		Query:                 query,
		Limit:                 10,
		AttributesToCrop:      []string{"message:200"},
		AttributesToHighlight: []string{"message"},
		Filters:               filter,
	})

	if err != nil {
		return err
	}

	for _, v := range res.Hits {
		isMap, ok := v.(map[string]interface{})
		formatted := isMap["_formatted"]
		mail := Mail{}
		if ok {
			if isFormatted, ok := formatted.(map[string]interface{}); ok {
				if isString, ok := isFormatted["message"].(string); ok {
					if strings.Contains(isString, "<em>") {
						isString = strings.Replace(isString, "<em>", yellow, -1)
						isString = strings.Replace(isString, "</em>", reset, -1)

						isMap["message"] = isString
						mail.Body = isString
					}
				}
			} else {
				mail.Body = get("message", isMap)
			}
			mail.From = get("from", isMap)
			mail.To = get("to", isMap)
			mail.Cc = get("cc", isMap)
			mail.Subject = get("subject", isMap)
			mail.Folder = get("folder", isMap)
			mail.Date = get("date", isMap)
		}
		println(mail.String())
		println("\n\n")
		println("===============================================")
	}

	return nil
}

func get(key string, container map[string]interface{}) string {
	val, ok := container[key].(string)
	if !ok {
		return ""
	}
	return val
}
