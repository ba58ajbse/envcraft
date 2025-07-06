package comment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeNewLines(t *testing.T) {
	tests := map[string]struct {
		orgLines []string
		line     int
		value    string
		want     []string
		wantErr  bool
	}{
		"append to end when line==0 and empty": {
			orgLines: []string{},
			line:     0,
			value:    "test comment",
			want:     []string{"# test comment"},
			wantErr:  false,
		},
		"append to end when line==0 and not empty, no trailing newline": {
			orgLines: []string{"FOO=\"bar\""},
			line:     0,
			value:    "test comment",
			want:     []string{"FOO=\"bar\"\n", "# test comment"},
			wantErr:  false,
		},
		"insert at line 1": {
			orgLines: []string{"FOO=\"bar\"\n", "BAR=\"baz\""},
			line:     1,
			value:    "test comment",
			want:     []string{"# test comment\n", "FOO=\"bar\"\n", "BAR=\"baz\""},
			wantErr:  false,
		},
		"insert at line greater than length": {
			orgLines: []string{"FOO=\"bar\""},
			line:     3,
			value:    "test comment",
			want:     []string{"FOO=\"bar\"\n", "\n", "# test comment"},
			wantErr:  false,
		},
		"empty lines, insert at line 2": {
			orgLines: []string{},
			line:     2,
			value:    "test comment",
			want:     []string{"\n", "# test comment"},
			wantErr:  false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &CommentCmd{
				Options: CommentOptions{
					Line:  tt.line,
					Value: tt.value,
				},
				OrgLines: tt.orgLines,
			}
			got, err := cmd.makeNewLines()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
