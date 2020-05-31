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
package indexer

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	filterAfter     = regexp.MustCompile(`after=([0-9-]+)`)
	filterBefore    = regexp.MustCompile(`before=([0-9-]+)`)
	filterTimeRange = regexp.MustCompile(`time=\"([0-9-]+):([0-9-]+)\"`)
)

// Filter is structured filter from user to meilisearch.
type Filter struct {
	query string

	After  time.Time
	Before time.Time
}

// NewFilter parses human-formatted filter into meilisearch filter.
func NewFilter(query string) *Filter {
	f := &Filter{
		query: query,
	}

	// find either time range or after & before
	if match := filterTimeRange.FindAllStringSubmatch(query, 2); len(match) > 0 {
		f.After = parseDate(match[0][1])
		f.Before = parseDate(match[0][2])
		out := fmt.Sprintf("date>%d AND date<%d", f.After.Unix(), f.Before.Unix())
		f.query = filterTimeRange.ReplaceAllString(f.query, out)
	} else {
		if match := filterAfter.FindStringSubmatch(query); len(match) > 0 {
			f.After = parseDate(match[1])
			out := fmt.Sprintf("date>%d", f.After.Unix())
			f.query = filterAfter.ReplaceAllString(f.query, out)
		}
		if match := filterBefore.FindStringSubmatch(query); len(match) > 0 {
			f.Before = parseDate(match[1])
			out := fmt.Sprintf("date<%d", f.Before.Unix())
			f.query = filterBefore.ReplaceAllString(f.query, out)
		}
	}
	return f
}

// parse datetime of format 2020-01-02 or just 2020
func parseDate(date string) time.Time {
	format := ""
	dashes := strings.Count(date, "-")
	if dashes == 0 {
		format = "2006"
	} else if dashes == 1 {
		format = "2006-01"
	} else if dashes == 2 {
		format = "2006-01-02"
	}

	ts, err := time.Parse(format, date)
	if err != nil {
		return time.Unix(0, 0)
	}
	return ts
}

func (f *Filter) Query() string {
	return f.query
}
