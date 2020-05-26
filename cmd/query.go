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
	"github.com/spf13/cobra"
	"strings"
	"tryffel.net/go/meilindex/indexer"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query mails with optional filters",
	Long: `Examples: 

1. 'meilindex query my mail' => match 'my mail'
2. 'meilindex query --folder inbox --subject "item received" my mail' => match 'my mail' in folder 'inbox' and subject '
item received'
`,
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().String("folder", "", "Folder to limit search to")
	queryCmd.Flags().String("from", "", "From sender (must match exactly)")
	queryCmd.Flags().String("to", "", "To receiver (must match exactly)")
	queryCmd.Flags().String("subject", "", "Subject (must match exactly)")

	queryCmd.Run = query
}

func query(cmd *cobra.Command, args []string) {
	q := strings.Join(args, " ")

	filter := ""

	folder, err := queryCmd.Flags().GetString("folder")
	if err == nil && folder != "" {
		filter += "folder=" + folder
	}

	indexer.SearchMail(q, filter)

}
