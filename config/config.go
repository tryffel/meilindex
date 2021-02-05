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

package config

const Version = "v0.1.0"

var Conf *Config

// Config is application configuration struct
type Config struct {
	File        File
	Imap        Imap
	Meilisearch Meilisearch
	Gui         Gui
}

// File is email locating on filesystem
type File struct {
	Directory string
	Recursive bool
	Mode      string
	BatchSize int
}

// Imap is imap-based email source
type Imap struct {
	Url              string
	Tls              bool
	SkipVerification bool
	Username         string
	Password         string
	Folder           string
}

// Meilisearch contains meilisearch-instance configuration
type Meilisearch struct {
	Url    string
	Index  string
	ApiKey string
}

type Gui struct {
	Mouse bool
}
