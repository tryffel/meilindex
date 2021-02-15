/*
 * Meilindex - mail indexing and search tool.
 * Copyright (C) 2021 Tero Vierimaa
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
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"tryffel.net/go/meilindex/config"
	"tryffel.net/go/meilindex/external"
)

func ReadVerbatimDir(path string, flushFunc func(mails []*Mail) error) error {
	var files []external.MboxFile
	var err error
	files, err = external.VerbatimFiles(path)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		logrus.Warning("Did not find any suitable Mbox files")
		return err
	}
	logrus.Infof("Found %d folders", len(files))

	batchSize := config.Conf.File.BatchSize

	mails := []*Mail{}

	totalMails := 0

	// read file and send email batches if batch is full. If smaller, return array of mails.
	// Try to always push reasonable batch size, even if single file contains less mails. Not batching
	// small files increases indexing time significantly.
	for i, v := range files {
		logrus.Infof("Indexing (%d / %d): %s", i, len(files), v.Name)
		mail, indexed, err := readVerbatimFile(v.File, v.Name, flushFunc)
		if err != nil {
			logrus.Error(err)
			continue
		}
		totalMails += indexed
		mails = append(mails, mail...)
		if len(mails) >= batchSize {
			err = flushFunc(mails)
			if err != nil {
				logrus.Error(err)
			}
			totalMails += len(mails)
			mails = []*Mail{}
		}
	}

	if len(mails) > 0 {
		err = flushFunc(mails)
		totalMails += len(mails)
	}
	logrus.Infof("Successfully indexed %d mails from %d files", totalMails, len(files))
	return nil
}

func readVerbatimFile(file, folder string, flushFunc func(mails []*Mail) error) ([]*Mail, int, error) {
	batchSize := 1000
	batch := 0
	currentBatchSize := 0
	totalMails := 0

	fd, err := os.Open(file)
	if err != nil {
		return nil, 0, err
	}

	msg, err := message.Read(fd)
	if err != nil {
		return nil, 0, fmt.Errorf("read file: %v", err)

	}

	var mails []*Mail

	if err != nil {
		if err == io.EOF {
			return nil, 0, nil
		}
		return mails, 0, err
	}

	if msg == nil {
		return nil, 0, nil
	}

	parsed, err := mail.CreateReader(fd)
	if err != nil {
		logrus.Warningf("(skip) parse mail: %v", err)
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

	err = fd.Close()
	if err != nil {
		logrus.Errorf("close file: %v", err)
	}

	if batch > 0 {
		logrus.Infof("Flushed %d batches, %d mails", batch, totalMails)
	}

	return mails, totalMails, nil
}
