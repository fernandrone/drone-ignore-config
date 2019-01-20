package ignorer

import (
	"bytes"
	"strings"
	"testing"
)

func TestIgnorer(t *testing.T) {
	var tests = []struct {
		ignore []string
		files  []string
		want   bool
	}{
		{[]string{""}, []string{""}, false},
		{[]string{""}, []string{"src/README.md"}, false},
		{[]string{"*"}, []string{""}, true},
		{[]string{"*"}, []string{"src/README.md"}, true},
		{[]string{"*", "!src/README.md"}, []string{"src/README.md"}, false},
		{[]string{"src/README.md"}, []string{"src/README.md", "src/DCOS.md"}, false},
		{[]string{"*.md"}, []string{"src/README.md", "src/DCOS.md"}, true},
		{[]string{"src"}, []string{"src/README.md", "src/DCOS.md"}, true},
		{[]string{"src/**"}, []string{"src/README.md", "src/DCOS.md"}, true},
		{[]string{"src/**"}, []string{"src/README.md", "src/DCOS.md", "data/file.txt"}, false},
		{[]string{"src/README.md", "src/DCOS.md"}, []string{"src/README.md", "src/DCOS.md"}, true},
		{[]string{"src/README.md", "src/DCOS.md"}, []string{"src/README.md"}, true},
		{[]string{"src/README.md", "src/DCOS.md"}, []string{"src/DCOS.md"}, true},
	}
	for _, test := range tests {
		in := bytes.NewReader([]byte(strings.Join(test.ignore, "\n")))
		di, err := New(in)

		if err != nil {
			t.Errorf("New(%q):\n%v", test.ignore, err)
		}

		if got := di.ShouldIgnore(test.files...); got != test.want {
			t.Errorf("New(%q).ShouldIgnore(%q):\nExpected \"%t\", got \"%t\"", test.ignore, test.files, test.want, got)
		}
	}
}
