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
			input:   "FOO=bar\nBAR=baz\n",
			want:    map[string]string{"FOO": "bar", "BAR": "baz"},
			wantErr: false,
		},
		"with comments": {
			input:   "# comment\nBAR= baz\n",
			want:    map[string]string{"BAR": "baz"},
			wantErr: false,
		},
		"with inline comment": {
			input:   "FOO=bar # inline\nBAR=baz #inline",
			want:    map[string]string{"FOO": "bar", "BAR": "baz"},
			wantErr: false,
		},
		"quoted value with inline comment": {
			input:   "FOO=\"bar\" # inline\nBAR='baz' #inline",
			want:    map[string]string{"FOO": "bar", "BAR": "baz"},
			wantErr: false,
		},
		"with spaces": {
			input:   "FOO = bar\nBAR= baz\nKEY =value",
			want:    map[string]string{"FOO": "bar", "BAR": "baz", "KEY": "value"},
			wantErr: false,
		},
		"quoted values": {
			input:   "FOO=\"bar\"\nBAR='baz'\n",
			want:    map[string]string{"FOO": "bar", "BAR": "baz"},
			wantErr: false,
		},
		"empty lines": {
			input:   "\n\nFOO=bar\n\nBAR=baz\n\n",
			want:    map[string]string{"FOO": "bar", "BAR": "baz"},
			wantErr: false,
		},
		"invalid line": {
			input:   "FOO\nBAR=baz\n",
			want:    nil,
			wantErr: true,
		},
		"value with # in string": {
			input:   "FOO=bar#notcomment\nBAR=baz #comment\nKEY=value# notcomment",
			want:    map[string]string{"FOO": "bar#notcomment", "BAR": "baz", "KEY": "value# notcomment"},
			wantErr: false,
		},
		"value with # in string and quoted": {
			input:   "FOO=\"bar # notcomment\"\nBAR='baz #notcomment'",
			want:    map[string]string{"FOO": "bar # notcomment", "BAR": "baz #notcomment"},
			wantErr: false,
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
