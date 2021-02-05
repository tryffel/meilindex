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

package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"tryffel.net/go/meilindex/config"
	"tryffel.net/go/meilindex/indexer"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index imap|file|dir",
	Short: "Index mails",
	Long: `Index mails from imap or file(s).
	
default imap / file / dir configuration is gathered from config file.
'dir' indexes config.file.directory 

Examples:
* meilindex index imap 
* meilindex index imap --folder Archive/All
* meilindex index file ~/.thunderbird/my-profile/ImapMail/host/Inbox
* meilindex index dir
* meilindex index dir ~/.thunderbird/my-profile/ImapMail/host
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("expected at least 1 argument")
		}
		if args[0] != "imap" && args[0] != "file" && args[0] != "dir" && args[0] != "mailspring" {
			return fmt.Errorf("expect location either 'imap', 'file', 'dir' or 'mailspring'")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	indexCmd.Flags().String("folder", "INBOX", "Imap folder to index")
	indexCmd.Run = indexMail
}

func indexMail(cmd *cobra.Command, args []string) {
	var mails []*indexer.Mail
	var err error
	meili, err := indexer.NewMeiliSearch()
	if err != nil {
		logrus.Errorf("Connect to meilisearch: %v", err)
		return
	}

	if args[0] == "file" {
		var file string
		var err error
		if len(args) >= 2 {
			file = args[1]
		} else {
			file, err = indexCmd.Flags().GetString("file")
		}
		mails, err = indexer.ReadFiles(file, false, meili.IndexMailBackground)
		if err != nil {
			logrus.Error(err)
		}
		meili.WaitIndexComplete()
	} else if args[0] == "dir" {
		recursive := true

		var dir string
		var err error
		if len(args) >= 2 {
			dir = args[1]
		} else {
			dir, err = indexCmd.Flags().GetString("dir")
		}

		if dir == "" {
			dir = viper.GetString("file.directory")
			recursive = viper.GetBool("file.recursive")
		}

		mails, err = indexer.ReadFiles(dir, recursive, meili.IndexMailBackground)
		if err != nil {
			fmt.Println(err)
			//return
		}
		meili.WaitIndexComplete()
	} else if args[0] == "mailspring" {
		var file string
		var err error
		if len(args) >= 2 {
			file = args[1]
		} else {
			file, err = indexCmd.Flags().GetString("file")
		}

		_, err = indexer.ReadMailspring(file, false, meili.IndexMailBackground)
		if err != nil {
			logrus.Error(err)
		}
		meili.WaitIndexComplete()
	} else {
		mails, err = retrieveImap()
		if err != nil {
			fmt.Printf("Error indexing from imap: %v\n", err)
			return
		}
		err = meili.IndexMail(mails)
		if err != nil {
			fmt.Printf("Error pushing mails to meilisearch: %v\n", err)
			return
		}
	}
}

func retrieveImap() ([]*indexer.Mail, error) {
	client := &indexer.Imap{
		Url:                 config.Conf.Imap.Url,
		Tls:                 config.Conf.Imap.Tls,
		TlsSkipVerification: config.Conf.Imap.SkipVerification,
		Username:            config.Conf.Imap.Username,
		Password:            config.Conf.Imap.Password,
	}
	var err error
	err = client.Connect()
	if err != nil {
		return nil, err
	}

	defer client.Disconnect()

	folder := config.Conf.Imap.Folder
	if f, err := indexCmd.Flags().GetString("folder"); f != "INBOX" && err != nil {
		folder = f
	}

	fmt.Printf("Index imap folder %s\n", folder)
	err = client.SelectMailbox(folder)
	if err != nil {
		return nil, fmt.Errorf("select folder: %v", err)
	}
	mails, err := client.FetchMail()
	return mails, err
}
