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
	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

type QueryInput struct {
	*cview.Form
	query     *cview.InputField
	filter    *cview.InputField
	queryFunc func(string, string)
}

func NewQueryInput(query func(query, filter string)) *QueryInput {
	q := &QueryInput{
		Form:      cview.NewForm(),
		query:     cview.NewInputField(),
		filter:    cview.NewInputField(),
		queryFunc: query,
	}

	q.SetBorder(false)
	q.query.SetLabel("Query")
	q.query.SetPlaceholder("marketing")

	q.filter.SetLabel("Filter")
	q.filter.SetPlaceholder("folder=inbox AND from=sender@mail.com")

	q.SetFieldTextColor(tcell.Color252)
	q.SetFieldBackgroundColor(tcell.Color235)

	q.AddFormItem(q.query)
	q.AddFormItem(q.filter)
	q.query.SetChangedFunc(q.search)
	return q
}

func (q *QueryInput) search(query string) {
	if q.queryFunc != nil {
		filter := q.filter.GetText()
		q.queryFunc(query, filter)
	}
}
