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
	"time"
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
		Limit:                 100,
		AttributesToHighlight: []string{"message", "subject"},
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
				body := getString("message", isFormatted)
				if body != "" {
					mail.Body = body
				}
				subject := getString("subject", isFormatted)
				if subject != "" {
					mail.Subject = subject
				}
			}
		} else {
			mail.Body = getString("message", isMap)
			mail.Subject = getString("subject", isMap)
		}
		mail.Uid = getString("uid", isMap)
		mail.Id = getString("id", isMap)
		mail.From = getString("from", isMap)
		mail.To = getStringArray("to", isMap)
		mail.Cc = getStringArray("cc", isMap)
		mail.Folder = getString("folder", isMap)
		mail.Timestamp = time.Unix(getInt("date", isMap), 0)
		mail.AttachmentNames = getStringArray("attachments", isMap)

		result[i] = mail

	}
	//println(mail.String())
	//println("\n\n")
	//p//rintln("===============================================")
	return result, int(res.ProcessingTimeMs), nil
}

func getString(key string, container map[string]interface{}) string {
	val, ok := container[key].(string)
	if !ok {
		return ""
	}
	return val
}

func getInt(key string, container map[string]interface{}) int64 {
	// int
	intVal, ok := container[key].(int)
	if ok {
		return int64(intVal)
	}
	int64Val, ok := container[key].(int64)
	if ok {
		return int64Val
	}
	float32Val, ok := container[key].(float32)
	if ok {
		return int64(float32Val)
	}
	float64Val, ok := container[key].(float64)
	if ok {
		return int64(float64Val)
	}
	return 0
}

func getStringArray(key string, container map[string]interface{}) []string {
	out := []string{}
	arr, ok := container[key].([]interface{})
	if !ok {
		return out
	}
	for _, v := range arr {
		if isString, ok := v.(string); ok {
			out = append(out, isString)
		}
	}
	return out
}
