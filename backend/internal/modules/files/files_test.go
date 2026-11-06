package files

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"easygym/internal/platform"
)

func TestSaveAndServe(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	png.Encode(&buf, image.NewRGBA(image.Rect(0, 0, 2, 2)))

	name, err := Save(dir, buf.Bytes())
	if err != nil || !nameRe.MatchString(name) || filepath.Ext(name) != ".png" {
		t.Fatalf("save: %q %v", name, err)
	}
	again, err := Save(dir, buf.Bytes())
	if err != nil || again != name {
		t.Fatalf("same bytes must give the same name: %q %v", again, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("want one stored file, have %d (temp files left behind?)", len(entries))
	}

	for _, bad := range [][]byte{[]byte("<svg onload=alert(1)>"), []byte("hello"), nil} {
		if _, err := Save(dir, bad); err == nil {
			t.Errorf("%q must be refused", bad)
		}
	}

	mux := http.NewServeMux()
	New(dir, &platform.Auth{}).Register(mux)
	get := func(path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		return rec
	}
	if rec := get("/api/files/" + name); rec.Code != 200 || rec.Header().Get("Content-Type") != "image/png" {
		t.Errorf("serve: %d %q", rec.Code, rec.Header().Get("Content-Type"))
	}
	if get("/api/files/..%2f..%2fetc%2fpasswd").Code != 404 {
		t.Error("only hash names are served")
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/api/admin/files", nil))
	if rec.Code != 401 {
		t.Errorf("upload without sign-in: %d", rec.Code)
	}
}
