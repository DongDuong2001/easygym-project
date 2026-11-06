package platform

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

// User is who signed in, from the verified token; Role comes from the accounts module.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"` // "customer", "staff" or "admin"
}

// Roles looks up (and on first sight records) a signed-in user's role; the accounts module implements it.
type Roles interface {
	Role(ctx context.Context, u User) (string, error)
}

// Auth checks the bearer token the web app sends after a Google sign-in. Any issuer that publishes
// a JWKS works, e.g. Supabase Auth at https://<ref>.supabase.co/auth/v1/.well-known/jwks.json.
// With no JWKS URL every protected route answers 401, so a missing setting never opens the admin.
type Auth struct {
	JWKSURL, Issuer, Audience string
	Roles                     Roles

	mu      sync.Mutex
	keys    map[string]crypto.PublicKey
	fetched time.Time
	now     func() time.Time // tests move the clock
}

// User returns the signed-in user with their role, or ErrUnauthorized when there is no valid token.
func (a *Auth) User(r *http.Request) (User, error) {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || a.JWKSURL == "" {
		return User{}, ErrUnauthorized
	}
	u, err := a.Verify(r.Context(), token)
	if err != nil {
		return User{}, err
	}
	u.Role, err = a.Roles.Role(r.Context(), u)
	return u, err
}

// Require runs h only for a signed-in user whose role is one of roles (any role when none are given).
func (a *Auth) Require(h func(http.ResponseWriter, *http.Request, User), roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := a.User(r)
		if err == nil && len(roles) > 0 && !slices.Contains(roles, u.Role) {
			err = ErrForbidden
		}
		if err != nil {
			WriteError(w, r, err)
			return
		}
		h(w, r, u)
	}
}

type claims struct {
	Sub      string          `json:"sub"`
	Email    string          `json:"email"`
	Iss      string          `json:"iss"`
	Aud      json.RawMessage `json:"aud"` // a string or a list of strings
	Exp      int64           `json:"exp"`
	Metadata struct {
		FullName string `json:"full_name"`
		Name     string `json:"name"`
	} `json:"user_metadata"`
}

// Verify checks an ES256 or RS256 JWT's signature, expiry, issuer and audience.
func (a *Auth) Verify(ctx context.Context, token string) (User, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return User{}, ErrUnauthorized
	}
	var head struct{ Alg, Kid string }
	if decodeSegment(parts[0], &head) != nil {
		return User{}, ErrUnauthorized
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return User{}, ErrUnauthorized
	}
	key, err := a.key(ctx, head.Kid)
	if err != nil {
		return User{}, err
	}
	digest := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	valid := false
	switch k := key.(type) {
	case *ecdsa.PublicKey:
		valid = head.Alg == "ES256" && len(sig) == 64 &&
			ecdsa.Verify(k, digest[:], new(big.Int).SetBytes(sig[:32]), new(big.Int).SetBytes(sig[32:]))
	case *rsa.PublicKey:
		valid = head.Alg == "RS256" && rsa.VerifyPKCS1v15(k, crypto.SHA256, digest[:], sig) == nil
	}
	var c claims
	if !valid || decodeSegment(parts[1], &c) != nil || c.Sub == "" {
		return User{}, ErrUnauthorized
	}
	now := time.Now
	if a.now != nil {
		now = a.now
	}
	if now().Unix() >= c.Exp || (a.Issuer != "" && c.Iss != a.Issuer) || (a.Audience != "" && !hasAudience(c.Aud, a.Audience)) {
		return User{}, ErrUnauthorized
	}
	name := c.Metadata.FullName
	if name == "" {
		name = c.Metadata.Name
	}
	return User{ID: c.Sub, Email: strings.ToLower(c.Email), Name: name}, nil
}

func decodeSegment(s string, v any) error {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func hasAudience(raw json.RawMessage, want string) bool {
	var one string
	if json.Unmarshal(raw, &one) == nil {
		return one == want
	}
	var many []string
	return json.Unmarshal(raw, &many) == nil && slices.Contains(many, want)
}

// key returns the public key for kid, refetching the JWKS when the issuer has rotated keys,
// at most once a minute so a flood of bad tokens cannot hammer the issuer.
func (a *Auth) key(ctx context.Context, kid string) (crypto.PublicKey, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if k, ok := a.keys[kid]; ok {
		return k, nil
	}
	if time.Since(a.fetched) < time.Minute {
		return nil, ErrUnauthorized
	}
	a.fetched = time.Now()
	keys, err := fetchJWKS(ctx, a.JWKSURL)
	if err != nil {
		return nil, err
	}
	a.keys = keys
	if k, ok := keys[kid]; ok {
		return k, nil
	}
	return nil, ErrUnauthorized
}

func fetchJWKS(ctx context.Context, url string) (map[string]crypto.PublicKey, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch jwks: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch jwks: status %d", res.StatusCode)
	}
	var set struct {
		Keys []struct{ Kid, Kty, Crv, X, Y, N, E string }
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&set); err != nil {
		return nil, fmt.Errorf("decode jwks: %w", err)
	}
	keys := map[string]crypto.PublicKey{}
	for _, k := range set.Keys {
		switch {
		case k.Kty == "EC" && k.Crv == "P-256":
			x, errX := base64.RawURLEncoding.DecodeString(k.X)
			y, errY := base64.RawURLEncoding.DecodeString(k.Y)
			if errX != nil || errY != nil || len(x) != 32 || len(y) != 32 {
				continue
			}
			pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), append(append([]byte{4}, x...), y...))
			if err == nil {
				keys[k.Kid] = pub
			}
		case k.Kty == "RSA":
			n, errN := base64.RawURLEncoding.DecodeString(k.N)
			e, errE := base64.RawURLEncoding.DecodeString(k.E)
			if errN != nil || errE != nil || len(e) > 4 {
				continue
			}
			keys[k.Kid] = &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(new(big.Int).SetBytes(e).Int64())}
		}
	}
	return keys, nil
}
