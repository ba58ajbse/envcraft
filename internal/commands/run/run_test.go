package run

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseEnv(t *testing.T) {
	tests := map[string]struct {
		input   string
		want    map[string]string
		wantErr bool
	}{
		"simple key-value": {
			input: "FOO=bar\nBAR=baz\n",
			want:  map[string]string{"FOO": "bar", "BAR": "baz"},
		},
		"with comments": {
			input: "# comment\nBAR= baz\n",
			want:  map[string]string{"BAR": "baz"},
		},
		"with inline comment": {
			input: "FOO=bar # inline\nBAR=baz #inline",
			want:  map[string]string{"FOO": "bar", "BAR": "baz"},
		},
		"quoted value with inline comment": {
			input: "FOO=\"bar\" # inline\nBAR='baz' #inline",
			want:  map[string]string{"FOO": "bar", "BAR": "baz"},
		},
		"with spaces": {
			input: "FOO = bar\nBAR= baz\nKEY =value",
			want:  map[string]string{"FOO": "bar", "BAR": "baz", "KEY": "value"},
		},
		"quoted values": {
			input: "FOO=\"bar\"\nBAR='baz'\n",
			want:  map[string]string{"FOO": "bar", "BAR": "baz"},
		},
		"empty lines": {
			input: "\n\nFOO=bar\n\nBAR=baz\n\n",
			want:  map[string]string{"FOO": "bar", "BAR": "baz"},
		},
		"invalid line": {
			input:   "FOO\nBAR=baz\n",
			want:    nil,
			wantErr: true,
		},
		"value with # in string": {
			input: "FOO=bar#notcomment\nBAR=baz #comment\nKEY=value# notcomment",
			want:  map[string]string{"FOO": "bar#notcomment", "BAR": "baz", "KEY": "value# notcomment"},
		},
		"value with # in string and quoted": {
			input: "FOO=\"bar # notcomment\"\nBAR='baz #notcomment'",
			want:  map[string]string{"FOO": "bar # notcomment", "BAR": "baz #notcomment"},
		},
		"single quotes containing double quotes": {
			input: `FOO='bar "hoge" fuga'`,
			want:  map[string]string{"FOO": "bar \"hoge\" fuga"},
		},
		"double quotes containing single quotes": {
			input: `FOO="bar 'hoge' fuga"`,
			want:  map[string]string{"FOO": "bar 'hoge' fuga"},
		},
		"single quotes containing backquotes": {
			input: "FOO='bar `hoge` fuga'",
			want:  map[string]string{"FOO": "bar `hoge` fuga"},
		},
		"double quotes containing backquotes": {
			input: "FOO=\"bar `hoge` fuga\"",
			want:  map[string]string{"FOO": "bar `hoge` fuga"},
		},
		"backquotes value": {
			input: "FOO=`bar`",
			want:  map[string]string{"FOO": "bar"},
		},
		"backquotes containing double quotes": {
			input: "FOO=`bar \"hoge\" fuga`",
			want:  map[string]string{"FOO": "bar \"hoge\" fuga"},
		},
		"backquotes containing single quotes": {
			input: "FOO=`bar 'hoge' fuga`",
			want:  map[string]string{"FOO": "bar 'hoge' fuga"},
		},
		"backquotes containing both double and single quotes": {
			input: "FOO=`bar \"hoge\" 'hoge' fuga`",
			want:  map[string]string{"FOO": "bar \"hoge\" 'hoge' fuga"},
		},
		"unquoted JSON value": {
			input: "FOO={\"hoge\": \"fuga\"}",
			want:  map[string]string{"FOO": "{\"hoge\": \"fuga\"}"},
		},
		"single quoted JSON value": {
			input: "FOO='{\"hoge\": \"fuga\"}'",
			want:  map[string]string{"FOO": "{\"hoge\": \"fuga\"}"},
		},
		"double quoted JSON value": {
			input: "FOO=\"{\"hoge\": \"fuga\"}\"",
			want:  map[string]string{"FOO": "{\"hoge\": \"fuga\"}"},
		},
		"backquoted JSON value": {
			input: "FOO=`{\"hoge\": \"fuga\"}`",
			want:  map[string]string{"FOO": "{\"hoge\": \"fuga\"}"},
		},
		"value contains =": {
			input: "FOO=bar==",
			want:  map[string]string{"FOO": "bar=="},
		},
		"unquoted value with newline": {
			input: "FOO=hoge\nfuga",
			want:  map[string]string{"FOO": "hoge\nfuga"},
		},
		"single quoted value with newline": {
			input: "FOO='hoge\nfuga'",
			want:  map[string]string{"FOO": "hoge\nfuga"},
		},
		"export statement": {
			input: "export Foo=hoge",
			want:  map[string]string{"Foo": "hoge"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := parseEnv(r)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := map[string]struct {
		args    []string
		wantErr bool
	}{
		"success: env command with .env": {
			args:    []string{"-f", "testdata/.env", "--", "env"},
			wantErr: false,
		},
		"error: missing file": {
			args:    []string{"-f", "testdata/notfound.env", "--", "env"},
			wantErr: true,
		},
		"error: missing command": {
			args:    []string{"-f", "testdata/.env"},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := Run(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
