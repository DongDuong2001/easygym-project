// Package files stores uploaded images on disk and serves them at /api/files/{name}.
// A file's name is its content hash, so uploading the same photo twice stores it once
// and a URL never changes meaning, which lets browsers cache it forever.
package files

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"easygym/internal/platform"
)

// MaxSize is the largest upload; photos from a phone fit, videos do not (those stay on YouTube).
const MaxSize = 10 << 20

// Only raster images: SVG can carry scripts, so it is refused.
var exts = map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/webp": ".webp", "image/gif": ".gif"}

var nameRe = regexp.MustCompile(`^[0-9a-f]{20}\.(jpg|png|webp|gif)$`)

type Module struct {
	dir  string
	auth *platform.Auth
}

func New(dir string, auth *platform.Auth) *Module { return &Module{dir: dir, auth: auth} }

func (m *Module) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/admin/files", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		r.Body = http.MaxBytesReader(w, r.Body, MaxSize+1<<20) // room for the multipart headers
		f, _, err := r.FormFile("file")
		if err != nil {
			platform.WriteError(w, r, platform.ValidationError{Field: "file", Message: "Chọn một ảnh tối đa 10 MB"})
			return
		}
		defer f.Close()
		data, err := io.ReadAll(io.LimitReader(f, MaxSize+1))
		if err == nil && len(data) > MaxSize {
			err = platform.ValidationError{Field: "file", Message: "Ảnh tối đa 10 MB"}
		}
		var name string
		if err == nil {
			name, err = Save(m.dir, data)
		}
		if err != nil {
			platform.WriteError(w, r, err)
			return
		}
		platform.WriteJSON(w, http.StatusCreated, map[string]string{"url": URL(name)})
	}, "staff", "admin"))
	mux.HandleFunc("GET /api/files/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if !nameRe.MatchString(name) {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.ServeFile(w, r, filepath.Join(m.dir, name))
	})
}

// URL is where a stored file is served.
func URL(name string) string { return "/api/files/" + name }

// Save stores data under its content hash and returns the file name. It refuses anything that is not
// a JPEG, PNG, WebP or GIF, judging by the bytes rather than the uploaded file name.
func Save(dir string, data []byte) (string, error) {
	ext, ok := exts[http.DetectContentType(data)]
	if !ok {
		return "", platform.ValidationError{Field: "file", Message: "Chỉ nhận ảnh JPG, PNG, WebP hoặc GIF"}
	}
	sum := sha256.Sum256(data)
	name := hex.EncodeToString(sum[:10]) + ext
	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); err == nil {
		return name, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	// Write then rename, so a reader never sees half a file.
	tmp, err := os.CreateTemp(dir, ".upload-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name()) // no-op after the rename
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	return name, os.Rename(tmp.Name(), path)
}
