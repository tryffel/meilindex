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

const (
	colorBackground         = tcell.Color234
	colorText               = tcell.Color252
	colorBackgroundSelected = tcell.Color23
	colorTextSelected       = tcell.Color252
	colorDisabled           = tcell.Color241
)

func init() {
	cview.Styles.PrimaryTextColor = colorText
	cview.Styles.PrimitiveBackgroundColor = colorBackground
}
