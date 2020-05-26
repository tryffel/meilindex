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
	"reflect"
	"testing"
)

func Test_stripdAddressNames(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    []string
	}{
		{
			address: `"Mickey Mouse" <mickey.mouse@gmail.com>`,
			want:    []string{"Mickey Mouse"},
		},
		{
			address: `"Person B" <person.b@gmail.com>, "Person C" <person.c@gmail.com>`,
			want:    []string{"Person B", "Person C"},
		},
		{
			address: "person.b@gmail.com, person.c@gmail.com",
			want:    []string{"person.b@gmail.com", "person.c@gmail.com"},
		},
		{
			address: `"'Person B'" <person.b@gmail.com>, "'Person C'" <person.c@gmail.com>`,
			want:    []string{"Person B", "Person C"},
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripdAddressNames(tt.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stripdAddressNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
