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
	query *cview.InputField
}

func NewQueryInput(query func(string)) *QueryInput {
	q := &QueryInput{
		Form:  cview.NewForm(),
		query: cview.NewInputField(),
	}

	q.SetBorder(true)
	q.SetTitle("Search")
	q.query.SetLabel("Query")
	q.query.SetPlaceholder("marketing OR sales")

	q.SetFieldTextColor(tcell.Color252)
	q.SetFieldBackgroundColor(tcell.Color235)

	q.AddFormItem(q.query)
	q.query.SetChangedFunc(query)
	return q
}
