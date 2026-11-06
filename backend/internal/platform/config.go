// Package platform is the shared plumbing every module uses: config, database pool and HTTP helpers.
// Modules depend on platform, never on each other's internals.
package platform

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	Addr        string
	// Sign-in tokens: JWKS URL, expected issuer and audience. Empty JWKS URL keeps every signed-in route closed.
	AuthJWKSURL, AuthIssuer, AuthAudience string
	AdminEmails                           []string // become admins on their first sign-in
	FilesDir                              string   // uploaded images
}

func LoadConfig() (Config, error) {
	loadDotEnv(".env", "../.env")
	c := Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		Addr:         os.Getenv("API_ADDR"),
		AuthJWKSURL:  os.Getenv("AUTH_JWKS_URL"),
		AuthIssuer:   os.Getenv("AUTH_ISSUER"),
		AuthAudience: os.Getenv("AUTH_AUDIENCE"),
		FilesDir:     os.Getenv("FILES_DIR"),
	}
	for _, e := range strings.Split(os.Getenv("ADMIN_EMAILS"), ",") {
		if e = strings.ToLower(strings.TrimSpace(e)); e != "" {
			c.AdminEmails = append(c.AdminEmails, e)
		}
	}
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.FilesDir == "" {
		c.FilesDir = "files"
	}
	if c.DatabaseURL == "" {
		return c, errors.New("DATABASE_URL is not set (see .env.example)")
	}
	return c, nil
}

// loadDotEnv reads KEY=VALUE lines from the first file that exists. Real environment
// variables win, so production sets them normally and never ships a .env file.
func loadDotEnv(paths ...string) {
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		for line := range strings.Lines(string(b)) {
			k, v, ok := strings.Cut(strings.TrimSpace(line), "=")
			if !ok || strings.HasPrefix(k, "#") {
				continue
			}
			if _, set := os.LookupEnv(k); !set {
				os.Setenv(k, strings.Trim(v, `"'`))
			}
		}
		return
	}
}
