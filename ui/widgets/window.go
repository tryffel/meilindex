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
	"tryffel.net/go/twidgets"
)

type Window struct {
	*twidgets.ModalLayout
	query    *QueryInput
	app      *cview.Application
	list     *MessageList
	preview  *cview.TextView
	navBar   *twidgets.NavBar
	help     *Help
	settings *Settings
	client   *indexer.Meilisearch

	mails []*indexer.Mail
}

func NewWindow() *Window {
	w := &Window{
		app:         cview.NewApplication(),
		ModalLayout: twidgets.NewModalLayout(),
		help:        NewHelp(),
		settings:    NewSettings(),
		client: &indexer.Meilisearch{
			Url:    config.Conf.Meilisearch.Url,
			Index:  config.Conf.Meilisearch.Index,
			ApiKey: config.Conf.Meilisearch.ApiKey,
		},
		preview: cview.NewTextView(),
	}

	colors := twidgets.NavBarColors{
		Background:            tcell.Color234,
		ButtonBackground:      tcell.Color234,
		ButtonBackgroundFocus: tcell.Color234,
		Text:                  tcell.Color252,
		TextFocus:             tcell.Color252,
		Shortcut:              tcell.Color214,
		ShortcutFocus:         tcell.Color214,
	}
	w.navBar = twidgets.NewNavBar(&colors, w.handleNavbar)
	w.navBar.AddButton(cview.NewButton("Help"), tcell.KeyF1)
	w.navBar.AddButton(cview.NewButton("Open mail"), tcell.KeyF2)
	w.navBar.AddButton(cview.NewButton("Settings"), tcell.KeyF3)

	w.query = NewQueryInput(w.search)
	w.list = NewMessageList(w.showMessage)

	err := w.client.Connect()
	if err != nil {
	}

	grid := w.ModalLayout.Grid()

	grid.SetRows(1, 5, -1, -1, -1, -1, -1, -1, 5, 1)
	grid.SetColumns(1, -1, -1, -1, -1, -1, -1, -1, -1, 1)
	grid.SetBorder(true)
	grid.SetTitle("Meilindex")

	grid.AddItem(w.navBar, 0, 0, 1, 10, 1, 15, false)
	grid.AddItem(w.query, 1, 0, 1, 10, 5, 15, true)
	grid.AddItem(w.list, 2, 0, 8, 6, 5, 15, false)
	grid.AddItem(w.preview, 2, 6, 8, 4, 5, 15, false)

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
	filt := indexer.NewFilter(filter).Query()
	mails, _, err := w.client.Query(text, filt)
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
	text += "From: " + mail.HighlightedFrom() + "\n"
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

	attachmentNames := ""
	if len(mail.AttachmentNames) > 0 {
		attachmentNames = fmt.Sprintf("}\nAttachments (%d): ", len(mail.AttachmentNames))
		for i, attachment := range mail.AttachmentNames {
			if i > 0 {
				attachmentNames += ", "
			}
			attachmentNames += attachment
		}
	}

	text += fmt.Sprintf(`Date: %s
Subject: %s%s
------------
	
%s`,
		mail.DateTime(), mail.HighlightedSubject(), attachmentNames, mail.HighlightedBody())
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

	if key == tcell.KeyF3 {
		if w.settings.isOpen || w.help.isOpen {
			return event
		} else {
			w.settings.isOpen = true
			w.AddDynamicModal(w.settings, twidgets.ModalSizeMedium)
			w.app.SetFocus(w.settings)
		}
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

	if key == tcell.KeyF1 {
		if w.help.isOpen || w.settings.isOpen {
			return event
		} else {
			w.help.isOpen = true
			w.AddDynamicModal(w.help, twidgets.ModalSizeMedium)
			w.app.SetFocus(w.help)
		}
	}

	if key == tcell.KeyEscape {
		if w.help.isOpen {
			w.help.isOpen = false
			w.RemoveModal(w.help)
			w.app.SetFocus(w.query)
		} else if w.settings.isOpen {
			w.settings.isOpen = false
			w.RemoveModal(w.settings)
			w.app.SetFocus(w.query)
		}
	}

	return event
}

func (w *Window) handleNavbar(label string) {

}

func init() {
	cview.Styles.PrimitiveBackgroundColor = tcell.Color234
	cview.Styles.PrimaryTextColor = tcell.Color252
	cview.Styles.BorderColor = tcell.Color246
	cview.Styles.TitleColor = tcell.Color252

}
