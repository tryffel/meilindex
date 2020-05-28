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
	"gitlab.com/tslocum/cview"
	"tryffel.net/go/meilindex/indexer"
	"tryffel.net/go/twidgets"
)

type MessageShort struct {
	*cview.TextView
	mail *indexer.Mail
}

func (m *MessageShort) SetSelected(selected twidgets.Selection) {
	switch selected {
	case twidgets.Selected:
		m.SetBackgroundColor(colorBackgroundSelected)
		m.SetTextColor(colorTextSelected)
	case twidgets.Blurred:
		m.SetBackgroundColor(colorDisabled)
	case twidgets.Deselected:
		m.SetBackgroundColor(colorBackground)
		m.SetTextColor(colorText)
	}
}

func NewMessageShort(index int, mail *indexer.Mail) *MessageShort {
	m := &MessageShort{
		TextView: cview.NewTextView(),
		mail:     mail,
	}

	m.SetBorder(false)
	m.SetDynamicColors(true)
	text := fmt.Sprintf(`%d. %s, %s
%s
`, index, mail.ShortDateTime(), mail.From, mail.HighlightedSubject())

	m.SetText(text)
	return m
}

type MessageList struct {
	*twidgets.ScrollList
	shortMessages []*MessageShort
	selectFunc    func(m *indexer.Mail)
}

func NewMessageList(selectFunc func(m *indexer.Mail)) *MessageList {
	m := &MessageList{

		selectFunc: selectFunc,
	}
	m.ScrollList = twidgets.NewScrollList(m.selectMail)
	m.ScrollList.Padding = 0
	m.SetBackgroundColor(colorBackground)

	return m
}

func (m *MessageList) AddMessage(index int, mail *indexer.Mail) {
	item := NewMessageShort(index, mail)
	m.AddItem(item)
	m.shortMessages = append(m.shortMessages, item)

}

func (m *MessageList) Clear() {
	m.ScrollList.Clear()
	m.shortMessages = []*MessageShort{}
}

func (m *MessageList) selectMail(index int) {
	if m.selectFunc != nil {
		mail := m.shortMessages[index].mail
		m.selectFunc(mail)
	}

}
