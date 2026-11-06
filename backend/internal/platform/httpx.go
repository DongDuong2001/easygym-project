package platform

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// Module is what each business module exposes to cmd/api.
type Module interface {
	Register(mux *http.ServeMux)
}

// ErrNotFound lets handlers map a missing row to 404 without knowing about pgx.
var ErrNotFound = errors.New("not found")

// ValidationError is a 422 with a message the frontend can show as is.
type ValidationError struct{ Field, Message string }

func (e ValidationError) Error() string { return e.Field + ": " + e.Message }

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false) // bodies are HTML from the CMS; escaping only bloats them
	if err := enc.Encode(v); err != nil {
		slog.Warn("write json", "err", err)
	}
}

// WriteError maps known errors to status codes and hides everything else behind a 500.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var ve ValidationError
	switch {
	case errors.As(err, &ve):
		WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": ve.Message, "field": ve.Field})
	case errors.Is(err, ErrNotFound):
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Không tìm thấy"})
	case errors.Is(err, ErrUnauthorized):
		WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Vui lòng đăng nhập"})
	case errors.Is(err, ErrForbidden):
		WriteJSON(w, http.StatusForbidden, map[string]string{"error": "Tài khoản không có quyền làm việc này"})
	default:
		slog.Error("request failed", "method", r.Method, "path", r.URL.Path, "err", err)
		WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Có lỗi xảy ra, vui lòng thử lại"})
	}
}

// Respond writes v as 200 JSON, or the error via WriteError.
func Respond(w http.ResponseWriter, r *http.Request, v any, err error) {
	if err != nil {
		WriteError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, v)
}

// DecodeJSON reads a small JSON body and rejects unknown fields. Requiring the JSON content type
// means a cross-site form cannot post here: browsers preflight JSON requests from other origins.
func DecodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	if ct, _, _ := strings.Cut(r.Header.Get("Content-Type"), ";"); ct != "application/json" {
		return ValidationError{Field: "body", Message: "Dữ liệu gửi lên không hợp lệ"}
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return ValidationError{Field: "body", Message: "Dữ liệu gửi lên không hợp lệ"}
	}
	return nil
}

var phoneRe = regexp.MustCompile(`^0\d{9}$`)

// NormalizePhone strips spaces, dots and dashes; ok is false unless it is 10 digits starting with 0.
// Same rule as isPhone in web/src/lib/shop.js.
func NormalizePhone(s string) (string, bool) {
	p := strings.NewReplacer(" ", "", ".", "", "-", "").Replace(s)
	return p, phoneRe.MatchString(p)
}

var slugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Slug checks a URL slug: lowercase letters, digits and single dashes, at most 120 characters.
func Slug(field, s string) error {
	if len(s) > 120 || !slugRe.MatchString(s) {
		return ValidationError{field, "Đường dẫn chỉ gồm chữ thường không dấu, số và dấu gạch ngang"}
	}
	return nil
}

// ImageURL checks an image address: an https URL or a file uploaded here (/api/files/...).
func ImageURL(field, s string, required bool) (string, error) {
	s, err := Text(field, "ảnh", s, required, 500)
	if err != nil || s == "" {
		return s, err
	}
	if !strings.HasPrefix(s, "https://") && !strings.HasPrefix(s, "/api/files/") {
		return s, ValidationError{field, "Ảnh phải là đường dẫn https:// hoặc ảnh đã tải lên"}
	}
	return s, nil
}

// Text trims s and checks its length in characters; label is the Vietnamese name shown in the error.
func Text(field, label, s string, required bool, max int) (string, error) {
	s = strings.TrimSpace(s)
	if required && s == "" {
		return "", ValidationError{field, "Vui lòng nhập " + label}
	}
	if utf8.RuneCountInString(s) > max {
		return "", ValidationError{field, fmt.Sprintf("Tối đa %d ký tự", max)}
	}
	return s, nil
}

// NewCode returns a prefix plus 8 random digits, e.g. EG04839201.
func NewCode(prefix string) string {
	b := make([]byte, 8)
	rand.Read(b)
	for i := range b {
		b[i] = '0' + b[i]%10
	}
	return prefix + string(b)
}

// Logging records method, path, status and duration for every request.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		defer func() {
			if v := recover(); v != nil {
				slog.Error("panic", "path", r.URL.Path, "panic", v)
				http.Error(rec, "internal error", http.StatusInternalServerError)
			}
			slog.Info("http", "method", r.Method, "path", r.URL.Path, "status", rec.status, "ms", time.Since(start).Milliseconds())
		}()
		next.ServeHTTP(rec, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
