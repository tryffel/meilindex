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
	"github.com/emersion/go-mbox"
	"github.com/emersion/go-message/mail"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"tryffel.net/go/meilindex/external"
)

// ReadFiles reads files and flushes batched mails to flushFunc
func ReadFiles(file string, recursive bool, flushFunc func(mails []*Mail) error) ([]*Mail, error) {
	var files []external.MboxFile
	var err error
	if recursive {
		files, err = external.MboxFiles(file, recursive)
		if err != nil {
			return nil, err
		}
		logrus.Infof("Indexing %d folders", len(files))
	} else {
		folder, _ := filepath.Abs(file)
		files = append(files, external.MboxFile{
			File: file,
			Name: folder,
		})
	}

	mails := []*Mail{}
	for _, v := range files {
		logrus.Infof("Index %s", v.Name)
		mail, err := readFile(v.File, v.Name, flushFunc)
		if err != nil {
			logrus.Error(err)
		} else {
			mails = append(mails, mail...)
		}
	}

	return mails, nil
}

func readFile(file, folder string, flushFunc func(mails []*Mail) error) ([]*Mail, error) {
	batchSize := 1000
	batch := 0
	currentBatchSize := 0
	totalMails := 0

	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	reader := mbox.NewReader(fd)

	msg, err := reader.NextMessage()
	var mails []*Mail

	for true {
		msg, err = reader.NextMessage()
		if err != nil {
			if err == io.EOF {
				break
			}
			return mails, err
		}

		if msg == nil {
			break
		}

		parsed, err := mail.CreateReader(msg)
		if err != nil {
			logrus.Errorf("parse mail: %v", err)
		} else {
			var m *Mail
			m, err = mailToMail(parsed)
			m.Folder = folder

			mails = append(mails, m)
			currentBatchSize += 1
			if currentBatchSize == batchSize {
				batch += 1
				totalMails += currentBatchSize
				err := flushFunc(mails)
				if err != nil {
					logrus.Errorf("Flush email batch: %v", err)
				}
				mails = []*Mail{}
				currentBatchSize = 0
			}
		}
	}
	if len(mails) > 0 {
		batch += 1
		totalMails += len(mails)
		err := flushFunc(mails)
		if err != nil {
			logrus.Errorf("Flush remaining email batch (size: %d): %v", len(mails), err)
		}
	}

	logrus.Infof("Flushed %d batches, %d mails", batch, totalMails)
	return mails, nil
}
