package platform

import (
	"net"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

// RateLimit lets each client IP POST to each of paths at most n times per window; other requests pass.
// The IP comes from X-Forwarded-For, which the SvelteKit proxy sets; the API must not be reachable
// any other way, or a client could send its own header.
// shortcut: counters live in this process, fine for one API instance; move them to the database if the API ever runs on several.
func RateLimit(next http.Handler, n int, window time.Duration, paths ...string) http.Handler {
	type slot struct {
		start time.Time
		count int
	}
	var mu sync.Mutex
	seen := map[string]*slot{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || !slices.Contains(paths, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		key := clientIP(r) + " " + r.URL.Path
		now := time.Now()
		mu.Lock()
		if len(seen) > 10_000 { // drop finished windows so the map cannot grow without end
			for k, s := range seen {
				if now.Sub(s.start) >= window {
					delete(seen, k)
				}
			}
		}
		s := seen[key]
		if s == nil || now.Sub(s.start) >= window {
			s = &slot{start: now}
			seen[key] = s
		}
		s.count++
		over := s.count > n
		mu.Unlock()
		if over {
			WriteJSON(w, http.StatusTooManyRequests, map[string]string{"error": "Bạn gửi quá nhiều lần, vui lòng thử lại sau ít phút"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		first, _, _ := strings.Cut(fwd, ",")
		return strings.TrimSpace(first)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
