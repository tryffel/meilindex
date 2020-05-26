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
	"tryffel.net/go/meilindex/config"
	"tryffel.net/go/meilindex/indexer"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index [location]",
	Short: "Index mails",
	Long: `Index mails from imap

Examples:
* meilindex index imap 
* meilindex index imap --folder INBOX
* meilindex index imap --folder Archive/Inbox
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("expected at least 1 argument")
		}
		if args[0] != "imap" && args[0] != "file" {
			return fmt.Errorf("expect location either 'imap' or 'file'")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)

	indexCmd.Flags().String("folder", "INBOX", "Folder to index")
	indexCmd.Flags().String("file", "", "File to index")
	indexCmd.Run = indexMail

}

func indexMail(cmd *cobra.Command, args []string) {
	var mails []*indexer.Mail
	var err error
	if args[0] == "file" {
		file, err := indexCmd.Flags().GetString("file")
		mails, err = indexer.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			//return
		}
	} else {
		mails, err = retrieveImap()
		if err != nil {
			fmt.Printf("Error indexing from imap: %v\n", err)
			return
		}
	}

	meili := indexer.Meilisearch{
		Url:    config.Conf.Meilisearch.Url,
		Index:  config.Conf.Meilisearch.Index,
		ApiKey: config.Conf.Meilisearch.ApiKey,
	}

	err = meili.Connect()
	if err != nil {
		fmt.Printf("Error connecting to meilisearch: %v\n", err)
		return
	}

	err = meili.IndexMail(mails)
	if err != nil {
		fmt.Printf("Error pushing mails to meilisearch: %v\n", err)
		return
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
