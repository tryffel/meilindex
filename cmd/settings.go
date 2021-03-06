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
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"tryffel.net/go/meilindex/indexer"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Configure indexing & ranking",
}

func init() {
	rootCmd.AddCommand(settingsCmd)
	settingsCmd.AddCommand(stopWordsCmd)
	settingsCmd.AddCommand(rankingCmd)
	settingsCmd.AddCommand(synonymsCmd)
	settingsCmd.PersistentFlags().Bool("get", true, "Use to get value")
	settingsCmd.PersistentFlags().Bool("set", false, "Use to set value")
	settingsCmd.Run = func(cmd *cobra.Command, args []string) {
		settingsCmd.Help()
	}

	stopWordsCmd.Run = stopWords
	rankingCmd.Run = rankings
	synonymsCmd.Run = synonyms
}

// stopWordsCmd represents the settings command
var stopWordsCmd = &cobra.Command{
	Use:   "stopwords get/set [file]",
	Short: "Configure stopwords",
}

// settingsCmd represents the settings command
var rankingCmd = &cobra.Command{
	Use:   "ranking get/set [file]",
	Short: "Configure ranking",
}

// settingsCmd represents the settings command
var synonymsCmd = &cobra.Command{
	Use:   "synonyms get/set [file]",
	Short: "Configure synonyms",
}

func stopWords(cmd *cobra.Command, args []string) {
	mode := "get"
	if len(args) > 1 {
		if args[0] == "set" {
			mode = "set"
		}
	}
	m, err := indexer.NewMeiliSearch()
	if err != nil {
		fmt.Printf("Error connecting to meilisearch: %v\n", err)
		return
	}
	if mode == "get" {
		stopWords, err := m.StopWords()
		if err != nil {
			fmt.Printf("Error getting stopwords: %v\n", err)
			return
		}

		fmt.Println(*stopWords)
	}
	if mode == "set" {
		file := args[1]
		fd, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer fd.Close()

		type Dto struct {
			StopWords []string `json:"stop_words"`
		}

		dto := Dto{}

		err = json.NewDecoder(fd).Decode(&dto)
		if err != nil {
			fmt.Printf("Error decoding json: %v\n", err)
			return
		}

		err = m.SetStopWords(dto.StopWords)
		if err != nil {
			fmt.Printf("Error applying stopwords: %v\n", err)
		}
	}
}

func rankings(cmd *cobra.Command, args []string) {
	mode := "get"
	if len(args) > 1 {
		if args[0] == "set" {
			mode = "set"
		}
	}
	m, err := indexer.NewMeiliSearch()
	if err != nil {
		fmt.Printf("Error connecting to meilisearch: %v\n", err)
		return
	}

	if mode == "get" {
		rules, err := m.RankingRules()
		if err != nil {
			fmt.Printf("Error getting ranking rules: %v\n", err)
			return
		}

		fmt.Println("Rules:")
		fmt.Println(*rules)
	} else {
		file := args[1]
		fd, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer fd.Close()

		type Dto struct {
			RankRules []string `json:"rankings"`
		}

		dto := Dto{}

		err = json.NewDecoder(fd).Decode(&dto)
		if err != nil {
			fmt.Printf("Error decoding json: %v\n", err)
			return
		}

		err = m.SetRankingRules(dto.RankRules)
		if err != nil {
			fmt.Printf("Error applying ranking rules: %v\n", err)
		}
	}
}

func synonyms(cmd *cobra.Command, args []string) {
	mode := "get"
	if len(args) > 1 {
		if args[0] == "set" {
			mode = "set"
		}
	}

	m, err := indexer.NewMeiliSearch()
	if err != nil {
		fmt.Printf("Error connecting to meilisearch: %v\n", err)
		return
	}

	if mode == "get" {
		synonyms, err := m.Synonyms()
		if err != nil {
			fmt.Printf("Error getting synonyms: %v\n", err)
			return
		}
		fmt.Println(*synonyms)
		return
	} else {
		file := args[1]
		fd, err := os.Open(file)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer fd.Close()

		type Dto struct {
			OneWay map[string][]string `json:"synonyms"`
		}

		dto := Dto{}

		err = json.NewDecoder(fd).Decode(&dto)
		if err != nil {
			fmt.Printf("Error decoding json: %v\n", err)
			return
		}

		if len(dto.OneWay) == 0 {
			logrus.Info("Clearing synonyms")
		}

		err = m.SetSynonyms(&dto.OneWay)
		if err != nil {
			logrus.Errorf("Update synonyms: %v", err)
		}
	}
}
