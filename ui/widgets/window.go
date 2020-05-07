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

type Window struct {
	*cview.Grid
	query   *QueryInput
	app     *cview.Application
	results *cview.Table
	preview *cview.TextView
	client  *indexer.Meilisearch

	mails []*indexer.Mail
}

func NewWindow() *Window {
	w := &Window{
		Grid:    cview.NewGrid(),
		app:     cview.NewApplication(),
		results: cview.NewTable(),
		client: &indexer.Meilisearch{
			Url:    config.MeilisearchHost,
			Index:  config.MeilisearchIndex,
			ApiKey: config.MeilisearchApiKey,
		},
		preview: cview.NewTextView(),
	}

	w.query = NewQueryInput(w.search)

	err := w.client.Connect()
	if err != nil {
	}
	w.SetRows(5, -1, -1)
	w.SetColumns(-1, 1)

	w.results.SetBorder(true)
	w.results.SetTitle("Results")
	w.results.SetSelectedStyle(tcell.Color252, tcell.Color23, 0)

	w.SetBorder(true)
	w.SetTitle("Meilindex")

	w.AddItem(w.query, 0, 0, 1, 1, 1, 15, true)
	w.AddItem(w.results, 1, 0, 1, 1, 5, 15, false)
	w.AddItem(w.preview, 2, 0, 1, 1, 5, 15, false)

	w.app.SetRoot(w, true).EnableMouse(true)
	w.app.SetFocus(w)

	w.results.SetSelectedFunc(w.showMessage)
	w.results.SetSelectable(true, false)
	w.preview.SetDynamicColors(true)
	w.preview.SetBorder(true)
	w.preview.SetTitle("Preview")
	w.app.SetInputCapture(w.inputCapture)
	return w
}

func (w *Window) Run() {
	w.app.Run()
}

func (w *Window) search(text string) {
	mails, err := w.client.Query(text, "")
	if err != nil {
		return
	}

	w.mails = mails
	w.results.Clear()
	w.results.SetCellSimple(0, 0, "#")
	w.results.SetCellSimple(0, 1, "From")
	w.results.SetCellSimple(0, 3, "Date")
	w.results.SetCellSimple(0, 4, "Subject")
	w.results.SetCellSimple(0, 5, "Message")
	for i, v := range mails {
		body := v.Body
		if strings.Contains(body, "<em>") {
			body = strings.Replace(body, "<em>", "[black:yellow:]", -1)
			body = strings.Replace(body, "</em>", "[-:-:-]", -1)
			v.Body = body
		}
		w.results.SetCellSimple(i+1, 0, fmt.Sprint(i+1))
		w.results.SetCellSimple(i+1, 1, v.From)
		w.results.SetCellSimple(i+1, 2, v.Date)
		w.results.SetCellSimple(i+1, 3, v.Subject)
		w.results.SetCellSimple(i+1, 4, v.Body)
	}
}

func (w *Window) showMessage(index int, col int) {
	if index == 0 {
		return
	}
	mail := w.mails[index-1]
	w.preview.SetText(mail.String())
}

func (w *Window) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	key := event.Key()
	if key == tcell.KeyTAB {
		var nextFocus cview.Primitive
		focused := w.app.GetFocus()

		switch focused {
		case w.results:
			nextFocus = w.preview
		case w.query, w.query.query:
			nextFocus = w.results
		case w.preview:
			nextFocus = w.query
		default:
			return event
		}

		w.app.SetFocus(nextFocus)
		return nil
	}
	return event
}

func init() {
	cview.Styles.PrimitiveBackgroundColor = tcell.Color234
	cview.Styles.PrimaryTextColor = tcell.Color252
	cview.Styles.BorderColor = tcell.Color246
	cview.Styles.TitleColor = tcell.Color252

}
