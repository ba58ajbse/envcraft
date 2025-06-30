package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MakeNewLines(t *testing.T) {
	tests := map[string]struct {
		orgLines []string
		l        int
		key      string
		value    string
		want     []string
	}{
		"append to end when l==0": {
			orgLines: []string{"FOO=\"bar\"\n", "BAR=\"baz\""},
			l:        0,
			key:      "NEW",
			value:    "value",
			want:     []string{"FOO=\"bar\"\n", "BAR=\"baz\"\n", "NEW=\"value\""},
		},
		"insert at line 1": {
			orgLines: []string{"FOO=\"bar\"\n", "BAR=\"baz\""},
			l:        1,
			key:      "NEW",
			value:    "value",
			want:     []string{"NEW=\"value\"\n", "FOO=\"bar\"\n", "BAR=\"baz\""},
		},
		"insert at line 2": {
			orgLines: []string{"FOO=\"bar\"\n", "BAR=\"baz\""},
			l:        2,
			key:      "NEW",
			value:    "value",
			want:     []string{"FOO=\"bar\"\n", "NEW=\"value\"\n", "BAR=\"baz\""},
		},
		"insert at line greater than length": {
			orgLines: []string{"FOO=\"bar\""},
			l:        4,
			key:      "NEW",
			value:    "value",
			want:     []string{"FOO=\"bar\"\n", "\n", "\n", "NEW=\"value\""},
		},
		"empty lines, append to end": {
			orgLines: []string{},
			l:        0,
			key:      "FOO",
			value:    "bar",
			want:     []string{"FOO=\"bar\""},
		},
		"empty lines, insert at line 1": {
			orgLines: []string{},
			l:        1,
			key:      "FOO",
			value:    "bar",
			want:     []string{"FOO=\"bar\""},
		},
		"empty lines, insert at line 2": {
			orgLines: []string{},
			l:        2,
			key:      "FOO",
			value:    "bar",
			want:     []string{"\n", "FOO=\"bar\""},
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			cmd := &AddCmd{
				Options: AddOptions{
					Line:  tt.l,
					Key:   tt.key,
					Value: tt.value,
				},
				OrgLines: tt.orgLines,
			}
			err := cmd.makeNewLines()
			if err != nil {
				t.Fatalf("makeNewLines() error = %v", err)
			}
			assert.Equal(t, tt.want, cmd.NewLines)
		})
	}
}

func Test_ParseAddOptions(t *testing.T) {
	tests := map[string]struct {
		opts    []string
		want    *AddOptions
		wantErr bool
	}{
		"normal case": {
			opts:    []string{"KEY", "VALUE", "-f", "test.env", "-l", "2"},
			want:    &AddOptions{Key: "KEY", Value: "VALUE", FilePath: "test.env", Line: 2},
			wantErr: false,
		},
		"missing key/value": {
			opts:    []string{"KEY"},
			want:    nil,
			wantErr: true,
		},
		"missing file": {
			opts:    []string{"KEY", "VALUE"},
			want:    nil,
			wantErr: true,
		},
		"negative line": {
			opts:    []string{"KEY", "VALUE", "-f", "test.env", "-l", "-1"},
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseAddOptions(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_readLines(t *testing.T) {
	tests := map[string]struct {
		filePath string
		want     []string
		wantErr  bool
	}{
		"normal": {
			filePath: "testdata/add_readLines_normal",
			want: []string{
				"# comment: This is a sample .env file\n",
				"ENV=\"dev\"\n",
				"\n",
				"DB_USER=\"user1\"\n",
				"DB_PASSWORD=\"pass123\"\n",
				"DB_HOST=\"localhost\"\n",
				"DB_NAME=\"sampledb\"\n",
				"DB_PORT=\"5432\"",
			},
			wantErr: false,
		},
		"empty file": {
			filePath: "testdata/add_readLines_empty",
			want:     []string{""},
			wantErr:  false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := &AddCmd{
				Options:  AddOptions{FilePath: tt.filePath},
				OrgLines: []string{},
			}
			err := cmd.readLines()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, cmd.OrgLines, "readLines() should return the correct lines")
			}
		})
	}
}

func Test_apply(t *testing.T) {
	tests := map[string]struct {
		newLines   []string
		wantOutput string
		wantErr    bool
	}{
		"normal": {
			newLines:   []string{"FOO=\"bar\"\n", "BAR=\"baz\""},
			wantOutput: "FOO=\"bar\"\nBAR=\"baz\"",
			wantErr:    false,
		},
		"empty": {
			newLines:   []string{},
			wantOutput: "",
			wantErr:    false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "envtest-*.env")
			assert.NoError(t, err)
			defer os.Remove(tmpfile.Name())
			tmpfile.Close()

			cmd := &AddCmd{
				Options:  AddOptions{FilePath: tmpfile.Name()},
				NewLines: tt.newLines,
			}
			err = cmd.apply()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				data, err := os.ReadFile(tmpfile.Name())
				assert.NoError(t, err)
				assert.Equal(t, tt.wantOutput, string(data))
			}
		})
	}
}
