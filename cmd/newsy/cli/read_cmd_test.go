package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	a1 = "\n[ID: 1]\nTitle: Title 1\nSource: The Tester\nTopics: tdd\n"
	a2 = "\n[ID: 2]\nTitle: Title 2\nSource: The Tester\nTopics: go\n"
	j1 = `{"id": 1, "title": "Title 1", "source": "The Tester", "topics": ["tdd"]}`
	j2 = `{"id": 2, "title": "Title 2", "source": "The Tester", "topics": ["go"]}`
)

func Test_ReadCmd(t *testing.T) {
	t.Parallel()

	pwd, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	mdb := fmt.Sprintf(`{"articles": [%s, %s], "topics": []}`, j1, j2)

	fp := filepath.Join(pwd, "newsy_db.json")
	err = ioutil.WriteFile(fp, []byte(mdb), 0644)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	noArgs := []string{}
	idList := []string{"1, 2"}

	testCases := []struct {
		name   string
		args   []string
		exp    string
		expErr bool
	}{
		{
			name:   "without any article ids",
			args:   noArgs,
			exp:    "please provide 1 or more article ids",
			expErr: true,
		}, {
			name:   "with only article ids",
			args:   idList,
			exp:    a1 + a2,
			expErr: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			bb := &bytes.Buffer{}
			cmd := &ReadCmd{}
			cmd.Out = bb

			ctx := context.Background()

			if err := cmd.Main(ctx, pwd, tC.args); err != nil {
				if !tC.expErr {
					t.Fatalf("expected err: %v, got err: %v", tC.expErr, err)
				}
				if tC.exp != err.Error() {
					t.Fatalf("expected: %v, got: %v", tC.exp, err.Error())
				}
				return
			}

			act := bb.String()

			if tC.exp != act {
				t.Fatalf("expected: %v, got: %v", tC.exp, act)
			}
		})
	}
}

// func createTestFS(t testing.TB) fstest.MapFS {
// 	t.Helper()
// 	a1 := `{"id": 1, "title": "Title 1", "source": "The Tester", "topics": ["tdd"]}`
// 	a2 := `{"id": 2, "title": "Title 2", "source": "The Tester", "topics": ["go"]}`
// 	a3 := `{"id": 3, "title": "Title 3", "source": "The Tester", "topics": ["neovim"]}`
//
// 	db := []byte(a1 + a2)
// 	db2 := []byte(a3)
//
// 	tfs := fstest.MapFS{
// 		"testdata":               {Mode: fs.ModeDir},
// 		"testdata/newsy_db.json": {Data: db},
// 		"testdata/alt_db.json":   {Data: db2},
// 	}
//
// 	return tfs
// }
