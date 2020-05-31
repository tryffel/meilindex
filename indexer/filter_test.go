package indexer

import (
	"testing"
)

func TestFilter_Query(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "after-date",
			query: `from=a AND to="b@mail.com" AND after=2020-01-02`,
			want:  `from=a AND to="b@mail.com" AND date>1577923200`,
		},
		{
			name:  "before-date",
			query: `from=a AND to="b@mail.com" AND before=2020-01-02`,
			want:  `from=a AND to="b@mail.com" AND date<1577923200`,
		},
		{
			name:  "after & before",
			query: `from=a AND to="b@mail.com" AND before=2020-01-02 AND after=2020-02-02`,
			want:  `from=a AND to="b@mail.com" AND date<1577923200 AND date>1580601600`,
		},
		{
			name:  "time-range",
			query: `from=a AND to="b@mail.com" AND time="2020-01-02:2020-02-02"`,
			want:  `from=a AND to="b@mail.com" AND date>1577923200 AND date<1580601600`,
		},
		{
			name:  "after-year",
			query: `from=a AND to="b@mail.com" AND after=2020`,
			want:  `from=a AND to="b@mail.com" AND date>1577836800`,
		},
		{
			name:  "after-month",
			query: `from=a AND to="b@mail.com" AND after=2020-01`,
			want:  `from=a AND to="b@mail.com" AND date>1577836800`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFilter(tt.query)
			if got := f.Query(); got != tt.want {
				t.Errorf("Query() = %v, want %v", got, tt.want)
			}
		})
	}
}
