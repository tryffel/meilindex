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
package widgets

import (
	"fmt"
	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
	"strings"
	"tryffel.net/go/meilindex/config"
	"tryffel.net/go/meilindex/indexer"
)

type Help struct {
	*cview.TextView

	infoText string

	page       int
	totalPages int

	isOpen bool
}

func (h *Help) SetDoneFunc(doneFunc func()) {
}

func (h *Help) SetVisible(visible bool) {
}

func NewHelp() *Help {
	h := &Help{
		TextView:   cview.NewTextView(),
		page:       0,
		totalPages: 4,
	}

	h.SetBackgroundColor(colorBackground)
	h.SetBorder(true)
	h.SetTitle("Help")
	h.SetBorderColor(tcell.Color230)
	h.SetTitleColor(colorText)
	h.SetDynamicColors(true)
	h.SetBorderPadding(0, 1, 2, 2)
	h.setContent()
	h.SetWordWrap(true)

	meili, err := indexer.NewMeiliSearch()
	if err == nil {
		stats := meili.Stats()

		h.infoText = fmt.Sprintf(`
[yellow]Meilisearch[-]:
Total mails: %d
Indexing in progress: %t
Server version: %s
`, stats.NumDocuments, stats.Indexing, stats.ServerVersion)
	}

	return h
}

func (h *Help) InputHandler() func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
		key := event.Key()
		if key == tcell.KeyLeft {
			if h.page > 0 {
				h.page -= 1
				h.setContent()
			}
		} else if key == tcell.KeyRight {
			if h.page < h.totalPages-1 {
				h.page += 1
				h.setContent()
			}
		} else {
			h.TextView.InputHandler()(event, setFocus)
		}
	}
}

func (h *Help) GetFocusable() cview.Focusable {
	return h.TextView.GetFocusable()
}

func (h *Help) setContent() {
	title := ""
	got := ""
	switch h.page {
	case 0:
		got = h.mainPage()
		title = "About"
	case 1:
		got = h.shortcutsPage()
		title = "Shortcuts"
	case 2:
		got = h.searchPage()
		title = "Searching"
	case 3:
		title = "Info"
		got = h.infoText
	default:
	}

	if title != "" {
		title = "[yellow::b]" + title + "[-::-]"
	}

	if got != "" {
		h.Clear()
		text := fmt.Sprintf("< %d / %d > %s \n\n", h.page+1, h.totalPages, title)
		text += got
		h.SetText(text)
		h.ScrollToBeginning()
	}
}

func (h *Help) mainPage() string {
	text := fmt.Sprintf("%s\n	[yellow]%s[-]\n\n", logo(), config.Version)
	text += "License: AGPL-v3, https://www.gnu.org/licenses/agpl-3.0.html"

	text += `
	
Email indexing and extremely fast full-text-search with Meilisearch. Meilindex supports configuring 
stop-words, ranking and synonyms. These are highly user-specific customizations and should be configured 
for more relevant search results. 

Features:
* Index mail from Imap or Mbox-file (tested with Thunderbird), store to Meilisearch
* Multiple configurations for different mailboxes
* Configure Meilisearch: stop words, ranking rules order
* Query Meilisearch instance either with CLI or with terminal gui
* Open selected mail in Thunderbird
	`
	return text
}

func (h *Help) shortcutsPage() string {
	return `[yellow]Movement[-]:
* Up/Down: Key up / down
* VIM-like keys: 
        * Up / Down: J / K, 
        * Top / Bottom of list: g / G 
        * Page Up / Down: Ctrl+F / Ctrl+B
* Switch between panels: Tab 
* Select button or item: Enter
* Close application: Ctrl-C
`
}

func (h *Help) searchPage() string {
	return `[yellow]Query[-]:
Query field supports full-text-search. Any field will be 
included, but only subject and message body will be highlighted. 
Query field must always include something for search results to appear, even with filters.
	
[yellow]Filter[-]:
You can define additional filters, which must match exactly. Boolean operators are supported. 
Supported fields are: [from, to, subject, cc, body, before/after/time]. 
	
Examples: 
	* 'folder=inbox AND from="example sender"'
	* 'folder=inbox AND NOT from="example.sender@example.company'
	
	
[yellow]Time range filters[-]:
Time ranges are parsed separately. 
Supported fields are: 'after', 'before' and 'time', see below for examples.
Format is year, optional month and optional day, 
e.g.: '2020', '2020-01' or '2020-01-01'. 2020 and 2020-01
will be expanded to 2020-01-01.

Examples:
	* 'after=2020-06 AND before=2021' (matches mails June 2020 - Jan 2021)
	* 'time="2020:2021"'
	
`
}

func logo() string {
	text := `
	 __  __      _ _ _           _
	|  \/  | ___(_) (_)_ __   __| | _____  __
	| |\/| |/ _ \ | | | '_ \ / _' |/ _ \ \/ /
	| |  | |  __/ | | | | | | (_| |  __/>  <
	|_|  |_|\___|_|_|_|_| |_|\__,_|\___/_/\_\
`
	return strings.TrimLeft(text, "\n")
}
