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
}

func LoadConfig() (Config, error) {
	loadDotEnv(".env", "../.env")
	c := Config{DatabaseURL: os.Getenv("DATABASE_URL"), Addr: os.Getenv("API_ADDR")}
	if c.Addr == "" {
		c.Addr = ":8080"
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
