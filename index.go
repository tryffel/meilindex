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
	"github.com/meilisearch/meilisearch-go"
	"github.com/sirupsen/logrus"
)

type Meilisearch struct {
	Url    string
	Index  string
	ApiKey string
	client *meilisearch.Client
}

func (m *Meilisearch) Connect() error {
	m.client = meilisearch.NewClient(meilisearch.Config{
		Host:   m.Url,
		APIKey: m.ApiKey,
	})

	indexExists := false
	_, err := m.client.Indexes().Get(m.Index)
	if err != nil {
		if e, ok := err.(*meilisearch.Error); ok {
			if e.StatusCode == 404 {
				err = nil
				indexExists = false
			}
		}

	} else {
		indexExists = true
	}

	if err != nil {
		return fmt.Errorf("get indexes")
	}

	if !indexExists {
		logrus.Warning("Creating new index")
		_, err = m.client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID:        m.Index,
			PrimaryKey: "id",
		})

	}

	if err != nil {
		return fmt.Errorf("create index: %v", err)
	}

	return nil
}

func (m *Meilisearch) IndexMail(mail []*Mail) error {

	logrus.Infof("Index %d mails", len(mail))

	documents := make([]map[string]interface{}, len(mail))

	for i, v := range mail {
		doc := map[string]interface{}{}
		doc["id"] = v.Id
		doc["date"] = v.Date
		doc["from"] = v.From
		doc["to"] = v.To
		doc["cc"] = v.Cc
		doc["subject"] = v.Subject
		doc["message"] = v.BodyPlainText()
		doc["folder"] = v.Folder
		documents[i] = doc
	}

	_, err := m.client.Documents(m.Index).AddOrUpdate(documents)

	if err != nil {
		return fmt.Errorf("push documents: %v", err)
	}

	logrus.Infof("Created / updated %d mails", len(mail))
	return nil
}
