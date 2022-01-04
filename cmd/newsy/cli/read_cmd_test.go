package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const (
	a1 = "\n[ID: 1]\nTitle: Title 1\nSource: The Tester\nTopics: tdd\n"
	a2 = "\n[ID: 2]\nTitle: Title 2\nSource: The Tester\nTopics: go\n"
	a3 = "\n[ID: 3]\nTitle: Title 3\nSource: The Tester\nTopics: neovim\n"
	a4 = "\n[ID: 4]\nTitle: Title 4\nSource: The Tester\nTopics: misc\n"
	j1 = `{"id": 1, "title": "Title 1", "source": "The Tester", "topics": ["tdd"]}`
	j2 = `{"id": 2, "title": "Title 2", "source": "The Tester", "topics": ["go"]}`
	j3 = `{"id": 3, "title": "Title 3", "source": "The Tester", "topics": ["neovim"]}`
	j4 = `{"id": 4, "title": "Title 4", "source": "The Tester", "topics": ["misc"]}`
)

func Test_ReadCmd(t *testing.T) {
	t.Parallel()

	pwd, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	mdb := fmt.Sprintf(`{"articles": [%s, %s, %s], "topics": []}`, j1, j2, j3)
	odb := fmt.Sprintf(`{"articles": [%s, %s], "topics": []}`, j1, j4)

	createTmpFile(t, pwd+"/newsy_db.json", mdb)
	createTmpFile(t, pwd+"/other_db.json", odb)
	createTmpFile(t, pwd+"/output_db.json", "")

	noArgs := []string{}
	idList := []string{"1", "2"}
	jsonIdList := []string{"-j", "1", "2"}
	fileIdList := []string{"-f", "other_db.json", "4"}
	outputIdList := []string{"-o", "output_db.json", "3"}

	testCases := []struct {
		name    string
		args    []string
		outFile string
		exp     string
		expErr  bool
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
		}, {
			name:   "with -j asJson flag and ids",
			args:   jsonIdList,
			exp:    mdb,
			expErr: false,
		}, {
			name:   "with -f flag and ids",
			args:   fileIdList,
			exp:    a4,
			expErr: false,
		}, {
			name:   "with -o location and ids",
			args:   outputIdList,
			exp:    a3,
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

			if cmd.Output == "" {
				act := bb.String()

				if tC.exp != act {
					t.Fatalf("expected: %v, got: %v", tC.exp, act)
				}
			} else {
				fp := pwd + "/" + cmd.Output
				assertFileEqual(t, fp, tC.exp)
			}
		})
	}
}

func assertFileEqual(t testing.TB, fp, exp string) {
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		t.Fatal(err)
	}

	act := string(b)

	if exp != act {
		t.Fatalf("expected: %v, got: %v", exp, act)
	}
}

func createTmpFile(t testing.TB, fp, data string) {
	t.Helper()

	if err := ioutil.WriteFile(fp, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(fp); errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
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
