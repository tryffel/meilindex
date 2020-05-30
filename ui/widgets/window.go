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
	"github.com/sirupsen/logrus"
	"gitlab.com/tslocum/cview"
	"tryffel.net/go/meilindex/config"
	"tryffel.net/go/meilindex/external"
	"tryffel.net/go/meilindex/indexer"
)

type Window struct {
	*cview.Grid
	query   *QueryInput
	app     *cview.Application
	list    *MessageList
	preview *cview.TextView
	client  *indexer.Meilisearch

	mails []*indexer.Mail
}

func NewWindow() *Window {
	w := &Window{
		Grid: cview.NewGrid(),
		app:  cview.NewApplication(),
		client: &indexer.Meilisearch{
			Url:    config.Conf.Meilisearch.Url,
			Index:  config.Conf.Meilisearch.Index,
			ApiKey: config.Conf.Meilisearch.ApiKey,
		},
		preview: cview.NewTextView(),
	}

	w.query = NewQueryInput(w.search)
	w.list = NewMessageList(w.showMessage)

	err := w.client.Connect()
	if err != nil {
	}
	w.SetRows(5, -1)
	w.SetColumns(-2, -1)

	w.SetBorder(true)
	w.SetTitle("Meilindex")

	w.AddItem(w.query, 0, 0, 1, 2, 1, 15, true)
	w.AddItem(w.list, 1, 0, 1, 1, 5, 15, false)
	w.AddItem(w.preview, 1, 1, 1, 1, 5, 15, false)

	w.app.SetRoot(w, true).EnableMouse(true)
	w.app.SetFocus(w)

	w.preview.SetDynamicColors(true)
	w.preview.SetBorder(true)
	w.preview.SetTitle("Preview")
	w.preview.SetWordWrap(true)
	w.app.SetInputCapture(w.inputCapture)
	return w
}

func (w *Window) Run() {
	w.app.Run()
}

func (w *Window) search(text, filter string) {
	mails, _, err := w.client.Query(text, filter)
	if err != nil {
		return
	}

	w.mails = mails

	w.list.Clear()
	for i, v := range mails {
		w.list.AddMessage(i+1, v)
	}
}

func (w *Window) showMessage(mail *indexer.Mail) {
	text := "Folder: " + mail.Folder + "\n"
	text += "From: " + mail.From + "\n"
	text += "To: "
	for i, v := range mail.To {
		if i > 0 {
			text += ", "
		}
		text += v
	}
	text += "\n"
	text += "Cc: "
	for i, v := range mail.Cc {
		if i > 0 {
			text += ", "
		}
		text += v
	}
	text += "\n"

	text += fmt.Sprintf(`Date: %s
Subject: %s
------------
	
%s`,
		mail.DateTime(), mail.HighlightedSubject(), mail.HighlightedBody())
	w.preview.SetText(text)
	w.preview.ScrollTo(0, 0)
}

func (w *Window) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	if key == tcell.KeyTAB {
		var nextFocus cview.Primitive
		focused := w.app.GetFocus()

		switch focused {
		case w.list:
			nextFocus = w.preview
		case w.query, w.query.query:
			nextFocus = w.query.filter
		case w.query.filter:
			nextFocus = w.list
		case w.preview:
			nextFocus = w.query
		default:
			return event
		}

		w.app.SetFocus(nextFocus)
		return nil
	}

	if key == tcell.KeyF2 {
		index := w.list.GetSelectedIndex()
		if index < len(w.list.shortMessages) {
			mail := w.list.shortMessages[index].mail
			err := external.OpenById(mail.Id)
			if err != nil {
				logrus.Errorf("Open mail in external application: %v", err)
			}
		}
	}

	return event
}

func init() {
	cview.Styles.PrimitiveBackgroundColor = tcell.Color234
	cview.Styles.PrimaryTextColor = tcell.Color252
	cview.Styles.BorderColor = tcell.Color246
	cview.Styles.TitleColor = tcell.Color252

}
