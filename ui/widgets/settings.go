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
	"tryffel.net/go/meilindex/indexer"
)

type Settings struct {
	*cview.TextView

	page       int
	totalPages int

	rankings  string
	stopwords string
	synonyms  string

	isOpen bool
}

func (s *Settings) SetDoneFunc(doneFunc func()) {
}

func (s *Settings) SetVisible(visible bool) {
}

func NewSettings() *Settings {
	s := &Settings{
		TextView:   cview.NewTextView(),
		page:       0,
		totalPages: 3,
	}

	s.SetBackgroundColor(colorBackground)
	s.SetBorder(true)
	s.SetTitle("Settings")
	s.SetBorderColor(tcell.Color230)
	s.SetTitleColor(colorText)
	s.SetDynamicColors(true)
	s.SetBorderPadding(0, 1, 2, 2)

	s.SetWordWrap(true)

	meili, err := indexer.NewMeiliSearch()
	if err == nil {
		rankings, err := meili.RankingRules()
		if err == nil {
			s.rankings = "- " + strings.Join(*rankings, "\n- ")
		}

		stopWords, err := meili.StopWords()
		if err == nil {
			s.stopwords = fmt.Sprintf("Total: %d\n\n", len(*stopWords))
			s.stopwords += strings.Join(*stopWords, ", ")
		}

		synonyms, err := meili.Synonyms()
		s.synonyms += fmt.Sprintf("Total: %d\n", len(*synonyms))
		if err == nil {
			for i, v := range *synonyms {
				s.synonyms += "\n- " + i + ": " + strings.Join(v, ", ")
			}
		}
	}
	s.setContent()
	return s
}

func (s *Settings) InputHandler() func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p cview.Primitive)) {
		key := event.Key()
		if key == tcell.KeyLeft {
			if s.page > 0 {
				s.page -= 1
				s.setContent()
			}
		} else if key == tcell.KeyRight {
			if s.page < s.totalPages-1 {
				s.page += 1
				s.setContent()
			}
		} else {
			s.TextView.InputHandler()(event, setFocus)
		}
	}
}

func (s *Settings) GetFocusable() cview.Focusable {
	return s.TextView.GetFocusable()
}

func (s *Settings) setContent() {
	title := ""
	got := ""
	switch s.page {
	case 0:
		title = "Ranking"
		got = s.rankings
	case 1:
		title = "Stop words"
		got = s.stopwords
	case 2:
		title = "Synonyms"
		got = s.synonyms
	default:
	}

	if title != "" {
		title = "[yellow::b]" + title + "[-::-]"
	}

	if got != "" {
		s.Clear()
		text := fmt.Sprintf("< %d / %d > %s \n\n", s.page+1, s.totalPages, title)
		text += got
		s.SetText(text)
		s.ScrollToBeginning()
	}
}
