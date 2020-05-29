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
	"github.com/spf13/cobra"
	"os"
	"strings"
	"tryffel.net/go/meilindex/config"
	"tryffel.net/go/meilindex/ui/widgets"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "meilindex",
	Short: "Email indexing and full-text-search with Meilisearch",
	Long: `Meilindex

Index emails from IMAP server or from local mail files. Running 'meilindex' opens gui for viewing indexed emails.
Licensed under AGPLv3`,
	Run: func(cmd *cobra.Command, args []string) {
		w := widgets.NewWindow()
		w.Run()

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.meilindex.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".meilindex" (without extension).
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		viper.SetConfigName(".meilindex")
	}

	viper.SetDefault("imap.url", "imap.mymail.com:993")
	viper.SetDefault("imap.tls", "true")
	viper.SetDefault("imap.skip_tls_verification", "false")
	viper.SetDefault("imap.username", "me@mymail.com")
	viper.SetDefault("imap.password", "memailing")
	viper.SetDefault("imap.folder", "INBOX")

	viper.SetDefault("file.directory", "/home/me/.mails")
	viper.SetDefault("file.recursive", "false")
	viper.SetDefault("file.mode", "thunderbird")

	viper.SetDefault("meilisearch.url", "http://localhost:7700")
	viper.SetDefault("meilisearch.index", "mail")
	viper.SetDefault("meilisearch.api_key", "masterKey")

	viper.SetEnvPrefix("meilindex")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}

	updateConfigFile()

	config.Conf = &config.Config{
		File: config.File{
			Directory: viper.GetString("file.directory"),
			Recursive: viper.GetBool("file.recursive"),
			Mode:      viper.GetString("file.mode"),
		},
		Imap: config.Imap{
			Url:              viper.GetString("imap.url"),
			Tls:              viper.GetBool("imap.tls"),
			SkipVerification: viper.GetBool("imap.skip_tls_verification"),
			Username:         viper.GetString("imap.username"),
			Password:         viper.GetString("imap.password"),
			Folder:           viper.GetString("imap.folder"),
		},
		Meilisearch: config.Meilisearch{
			Url:    viper.GetString("meilisearch.url"),
			Index:  viper.GetString("meilisearch.index"),
			ApiKey: viper.GetString("meilisearch.api_key"),
		},
	}

	viper.ConfigFileUsed()
}

func updateConfigFile() {

	err := viper.WriteConfig()
	if err == nil {
		return
	}
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		fmt.Println("Writing new config file")
		fd, err := os.Create(viper.ConfigFileUsed())
		if err != nil {
			fmt.Println(err)
		} else {
			fd.Close()
		}

		err = viper.SafeWriteConfig()
		if err == nil {
			return
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}
