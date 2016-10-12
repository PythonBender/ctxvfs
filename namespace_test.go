package ctxvfs

import (
	"testing"
	"time"
)

func TestNameSpace(t *testing.T) {
	// has mount point
	t1 := NameSpace{}
	t1.Bind("/fs1", Map(map[string][]byte{"fs1file": []byte("abcdefgh")}), "/", BindReplace)

	// has no mount point
	var t2 NameSpace

	testcases := map[string][]bool{
		"/":            []bool{true, false},
		"/fs1":         []bool{true, false},
		"/fs1/fs1file": []bool{true, false},
	}

	fss := []FileSystem{t1, t2}

	for j, fs := range fss {
		for k, v := range testcases {
			_, err := fs.Stat(nil, k)
			result := err == nil
			if result != v[j] {
				t.Errorf("fs: %d, testcase: %s, want: %v, got: %v, err: %s", j, k, v[j], result, err)
			}
		}
	}

	fi, err := t1.Stat(nil, "/")
	if err != nil {
		t.Fatal(err)
	}
	if fi.Name() != "/" {
		t.Errorf("t2.Name() : want:%s got:%s", "/", fi.Name())
	}
	if !fi.ModTime().IsZero() {
		t.Errorf("t2.Modime() : want:%v got:%v", time.Time{}, fi.ModTime())
	}
}

func TestNameSpace_ancestorDirs(t *testing.T) {
	mfs := Map(map[string][]byte{"a/b.txt": []byte("c")})
	fs := NameSpace{}
	fs.Bind("/x/y", mfs, "/", BindBefore)

	statTests := []struct {
		path      string
		wantIsDir bool
	}{
		{"/", true},
		{"/x", true},
		{"/x/y", true},
		{"/x/y/a", true},
	}
	for _, test := range statTests {
		fi, err := fs.Stat(nil, test.path)
		if err != nil {
			t.Errorf("Stat(%q): %s", test.path, err)
			continue
		}
		if fi.Mode().IsDir() != test.wantIsDir {
			t.Errorf("Stat(%q): got IsDir %v, want %v", test.path, fi.Mode().IsDir(), test.wantIsDir)
		}
	}
}