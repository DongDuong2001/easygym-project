package platform

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fixedRoles string

func (r fixedRoles) Role(context.Context, User) (string, error) { return string(r), nil }

func sign(t *testing.T, key *ecdsa.PrivateKey, kid string, claims map[string]any) string {
	t.Helper()
	enc := func(v any) string {
		b, _ := json.Marshal(v)
		return base64.RawURLEncoding.EncodeToString(b)
	}
	signed := enc(map[string]string{"alg": "ES256", "kid": kid}) + "." + enc(claims)
	digest := sha256.Sum256([]byte(signed))
	r, s, err := ecdsa.Sign(rand.Reader, key, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	sig := make([]byte, 64)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:])
	return signed + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func TestAuth(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := key.PublicKey.Bytes() // 0x04 || X || Y
	jwks := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{{
			"kid": "k1", "kty": "EC", "crv": "P-256",
			"x": base64.RawURLEncoding.EncodeToString(pub[1:33]),
			"y": base64.RawURLEncoding.EncodeToString(pub[33:]),
		}}})
	}))
	defer jwks.Close()

	now := time.Unix(1_800_000_000, 0)
	a := &Auth{JWKSURL: jwks.URL, Issuer: "https://x.supabase.co/auth/v1", Audience: "authenticated",
		Roles: fixedRoles("staff"), now: func() time.Time { return now }}
	good := map[string]any{"sub": "u1", "email": "An@Mail.com", "iss": "https://x.supabase.co/auth/v1",
		"aud": "authenticated", "exp": now.Unix() + 60, "user_metadata": map[string]string{"full_name": "An"}}
	with := func(k string, v any) map[string]any {
		c := map[string]any{}
		for kk, vv := range good {
			c[kk] = vv
		}
		c[k] = v
		return c
	}
	call := func(token string, roles ...string) int {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/me", nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		a.Require(func(w http.ResponseWriter, _ *http.Request, u User) {
			if u.ID != "u1" || u.Email != "an@mail.com" || u.Name != "An" || u.Role != "staff" {
				t.Errorf("user %+v", u)
			}
		}, roles...)(rec, req)
		return rec.Code
	}

	token := sign(t, key, "k1", good)
	other, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	cases := map[string]struct {
		token string
		roles []string
		want  int
	}{
		"valid, any role":    {token, nil, 200},
		"valid, staff route": {token, []string{"staff", "admin"}, 200},
		"wrong role":         {token, []string{"admin"}, 403},
		"no token":           {"", nil, 401},
		"garbage":            {"a.b.c", nil, 401},
		"expired":            {sign(t, key, "k1", with("exp", now.Unix())), nil, 401},
		"other issuer":       {sign(t, key, "k1", with("iss", "https://evil")), nil, 401},
		"audience list":      {sign(t, key, "k1", with("aud", []string{"x", "authenticated"})), nil, 200},
		"wrong audience":     {sign(t, key, "k1", with("aud", "anon")), nil, 401},
		"other key":          {sign(t, other, "k1", good), nil, 401},
		"unknown kid":        {sign(t, key, "k2", good), nil, 401},
	}
	for name, c := range cases {
		if got := call(c.token, c.roles...); got != c.want {
			t.Errorf("%s: status %d, want %d", name, got, c.want)
		}
	}

	if _, err := (&Auth{}).User(httptest.NewRequest("GET", "/", nil)); !errors.Is(err, ErrUnauthorized) {
		t.Errorf("unconfigured auth must refuse, got %v", err)
	}
}
