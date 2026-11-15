package webserver

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_getMediaFilterInfoFromUrlQuery(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
		want   MediaFilterInfo
	}{
		{
			name:   "empty query means no filter",
			values: url.Values{},
			want:   MediaFilterInfo{},
		},
		{
			name: "valid parse from url query",
			values: url.Values{
				"name":    {"  Matrix "},
				"minSize": {"1.5"},
				"okOnly":  {"true"},
			},
			want: MediaFilterInfo{Name: "Matrix", MinSize: 1.5, OkOnly: true},
		},
		{
			name:   "invalid values fall back to no filter",
			values: url.Values{"minSize": {"-3"}, "okOnly": {"yes"}},
			want:   MediaFilterInfo{},
		},
		{
			name:   "infinite size is ignored",
			values: url.Values{"minSize": {"Inf"}},
			want:   MediaFilterInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMediaFilterInfoFromUrlQuery(tt.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMediaFilterInfoFromUrlQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_MediaFilterInfo_FuzzyNameMatch(t *testing.T) {
	tests := []struct {
		name  string
		query string
		title string
		want  bool
	}{
		{name: "empty query matches everything", query: "", title: "Anything", want: true},
		{name: "dot, underscore and space are interchangeable", query: "2 Test", title: "2.Test", want: true},
		{name: "underscore title matches space query", query: "2 test", title: "2_Test", want: true},
		{name: "space title matches dot query", query: "2.test", title: "2 Test", want: true},
		{name: "dash matches space", query: "a b", title: "A-B", want: true},
		{name: "en dash matches space", query: "a b", title: "A–B", want: true},
		{name: "dash matches underscore", query: "a_b", title: "a-b", want: true},
		{name: "runs of separators collapse", query: "2 test", title: "2 . _ Test", want: true},
		{name: "round brackets are removed", query: "test 2020", title: "Test (2020)", want: true},
		{name: "square brackets are removed", query: "test 2020", title: "Test [2020]", want: true},
		{name: "curly and angle brackets are removed", query: "test2020x", title: "Test{2020}<x>", want: true},
		{name: "full-width brackets are removed", query: "test2020", title: "Test（2020）", want: true},
		{name: "colon, semicolon and exclamation are removed", query: "re test", title: "Re: Test;!", want: true},
		{name: "question mark and quotes are removed", query: "whos there", title: "Who's \"There\"?", want: true},
		{name: "removed characters do not separate words", query: "whos", title: "Who's", want: true},
		{name: "case is ignored", query: "MATRIX", title: "The matrix", want: true},
		{name: "partial match", query: "trix", title: "The Matrix", want: true},
		{name: "different characters do not match", query: "2 test", title: "3 Test", want: false},
		{name: "dot is not removed", query: "2 test", title: "2test", want: false},
		{name: "missing text does not match", query: "matrixx", title: "The Matrix", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NameFilter(tt.query)
			if got := filter.FuzzyNameMatch(tt.title); got != tt.want {
				t.Errorf("FuzzyNameMatch(%q) with query %q = %v, want %v", tt.title, tt.query, got, tt.want)
			}
		})
	}
}

func Test_SizeGB_Bytes(t *testing.T) {
	if got := SizeGB(0).Bytes(); got != 0 {
		t.Errorf("no size filter: Bytes() = %d, want 0", got)
	}
	if got := SizeGB(2).Bytes(); got != 2_000_000_000 {
		t.Errorf("2 GB: Bytes() = %d, want 2000000000", got)
	}
	if got := SizeGB(1.5).Bytes(); got != 1_500_000_000 {
		t.Errorf("1.5 GB: Bytes() = %d, want 1500000000", got)
	}
}
