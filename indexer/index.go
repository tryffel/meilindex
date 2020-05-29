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
	"crypto/md5"
	"fmt"
	"github.com/meilisearch/meilisearch-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"tryffel.net/go/meilindex/config"
)

// NewMeilisearch creates new connection.
func NewMeiliSearch() (*Meilisearch, error) {
	m := &Meilisearch{
		Url:    config.Conf.Meilisearch.Url,
		Index:  config.Conf.Meilisearch.Index,
		ApiKey: config.Conf.Meilisearch.ApiKey,
	}
	err := m.Connect()
	return m, err
}

// Meilisearch is a connector to Meilisearch.
type Meilisearch struct {
	Url    string
	Index  string
	ApiKey string
	client *meilisearch.Client
}

// Connect creates a connection to meilisearch instance and initializes index if neccessary.
func (m *Meilisearch) Connect() error {
	m.client = meilisearch.NewClient(meilisearch.Config{
		Host:   m.Url,
		APIKey: m.ApiKey,
	})

	m.client = meilisearch.NewClientWithCustomHTTPClient(meilisearch.Config{
		Host:   m.Url,
		APIKey: m.ApiKey,
	}, http.Client{
		Timeout: 10 * time.Second,
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
		return fmt.Errorf("get indexes: %v", err)
	}

	if !indexExists {
		logrus.Warning("Creating new index")
		_, err = m.client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID:        m.Index,
			PrimaryKey: "uid",
		})

	}

	if err != nil {
		return fmt.Errorf("create index: %v", err)
	}

	return nil
}

// IndexMail indexes new mail or updates existing mails.
func (m *Meilisearch) IndexMail(mail []*Mail) error {

	logrus.Infof("Index %d mails", len(mail))

	documents := make([]map[string]interface{}, len(mail))

	for i, v := range mail {
		v.Sanitize()
		doc := map[string]interface{}{}
		doc["id"] = v.Id
		doc["date"] = v.Timestamp.Unix()
		doc["from"] = v.From
		doc["to"] = v.To
		doc["cc"] = v.Cc
		doc["subject"] = v.Subject
		doc["message"] = v.Body
		doc["folder"] = v.Folder
		documents[i] = doc

		// email ids can be too complex for meilisearch. Use md5 as a unique id for mail.
		hash := md5.Sum([]byte(v.Uid))
		doc["uid"] = fmt.Sprintf("%x", hash)
	}

	res, err := m.client.Documents(m.Index).AddOrReplace(documents)

	logrus.Info("Meilisearch update id: ", res.UpdateID)

	if err != nil {
		return fmt.Errorf("push documents: %v", err)
	} else {
		logrus.Infof("Created / updated %d mails", len(mail))
	}
	return nil
}

// RankingRules returns a list of ranking rules. First rule is the most important, last is least important.
func (m *Meilisearch) RankingRules() (*[]string, error) {

	return m.client.Settings(m.Index).GetRankingRules()
}

func (m *Meilisearch) SetRankingRules(rules []string) error {
	_, err := m.client.Settings(m.Index).UpdateRankingRules(rules)
	return err
}

// StopWords returns all stop words currently being used.
func (m *Meilisearch) StopWords() (*[]string, error) {
	words, err := m.client.Settings(m.Index).GetStopWords()
	if err != nil {
		return nil, err
	}
	return words, nil
}

func (m *Meilisearch) SetStopWords(words []string) error {
	_, err := m.client.Settings(m.Index).UpdateStopWords(words)
	return err
}
