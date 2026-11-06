// Command mirror copies every image the database still loads from easygym.vn into FILES_DIR and
// points the database at the copies (/api/files/...). easygym.vn is run by another company and can
// disappear at any time. Running it again only picks up URLs not copied yet.
//
//	go run ./cmd/mirror            copy and rewrite
//	go run ./cmd/mirror -dry-run   only list what would be copied
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"easygym/internal/modules/files"
	"easygym/internal/platform"
)

// urlRe finds easygym.vn image URLs in plain columns, arrays and HTML bodies alike. Commas and braces
// end a URL because Postgres prints arrays as {url1,url2} without quotes.
var urlRe = regexp.MustCompile(`(?i)https?://(?:www\.)?easygym\.vn/[^\s"'<>(),{}]+\.(?:jpe?g|png|webp|gif)\b`)

// Every column that can hold an image URL. This command is the one place outside the modules
// that touches their tables, like cmd/seed.
var columns = []struct {
	table, col string
	array      bool
}{
	{"catalog.products", "image", false}, {"catalog.products", "images", true},
	{"catalog.products", "description", false}, {"catalog.products", "detail", false},
	{"catalog.product_questions", "answer", false},
	{"catalog.courses", "card", false}, {"catalog.courses", "image", false},
	{"catalog.courses", "summary", false}, {"catalog.courses", "body", false},
	{"content.posts", "image", false}, {"content.posts", "excerpt", false}, {"content.posts", "body", false},
	{"content.album_images", "url", false}, {"content.pages", "body", false},
}

func main() {
	dry := flag.Bool("dry-run", false, "only list the URLs that would be copied")
	flag.Parse()
	if err := run(*dry); err != nil {
		fmt.Fprintln(os.Stderr, "mirror:", err)
		os.Exit(1)
	}
}

func run(dry bool) error {
	cfg, err := platform.LoadConfig()
	if err != nil {
		return err
	}
	ctx := context.Background()
	db, err := platform.OpenDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	var urls []string
	for _, c := range columns {
		rows, _ := db.Query(ctx, fmt.Sprintf(`select %s::text from %s`, c.col, c.table))
		texts, err := pgx.CollectRows(rows, pgx.RowTo[string])
		if err != nil {
			return fmt.Errorf("read %s.%s: %w", c.table, c.col, err)
		}
		for _, t := range texts {
			urls = append(urls, FindURLs(t)...)
		}
	}
	slices.Sort(urls)
	urls = slices.Compact(urls)
	fmt.Printf("%d image URLs on easygym.vn\n", len(urls))
	if dry {
		for _, u := range urls {
			fmt.Println(u)
		}
		return nil
	}

	client := &http.Client{Timeout: 30 * time.Second}
	moved := map[string]string{}
	failed := 0
	for i, u := range urls {
		name, err := download(client, cfg.FilesDir, u)
		if err != nil {
			failed++
			fmt.Printf("[%d/%d] FAILED %s: %v\n", i+1, len(urls), u, err)
			continue
		}
		moved[u] = files.URL(name)
		fmt.Printf("[%d/%d] %s -> %s\n", i+1, len(urls), u, moved[u])
	}

	// Rewrite only after every download, in one transaction: the site never points at a missing file.
	err = platform.InTx(ctx, db, func(tx pgx.Tx) error {
		batch := &pgx.Batch{}
		for old, now := range moved {
			for _, c := range columns {
				if c.array {
					batch.Queue(fmt.Sprintf(`update %s set %s = array_replace(%[2]s, $1, $2) where $1 = any(%[2]s)`, c.table, c.col), old, now)
				} else {
					batch.Queue(fmt.Sprintf(`update %s set %s = replace(%[2]s, $1, $2) where strpos(%[2]s, $1) > 0`, c.table, c.col), old, now)
				}
			}
		}
		return tx.SendBatch(ctx, batch).Close()
	})
	if err != nil {
		return fmt.Errorf("rewrite urls: %w", err)
	}
	fmt.Printf("copied %d, failed %d; files are in %s\n", len(moved), failed, cfg.FilesDir)
	if failed > 0 {
		return fmt.Errorf("%d downloads failed, run again to retry them", failed)
	}
	return nil
}

// FindURLs returns the easygym.vn image URLs in s.
func FindURLs(s string) []string { return urlRe.FindAllString(s, -1) }

func download(client *http.Client, dir, url string) (string, error) {
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", res.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(res.Body, files.MaxSize+1))
	if err != nil {
		return "", err
	}
	if len(data) > files.MaxSize {
		return "", fmt.Errorf("larger than %d MB", files.MaxSize>>20)
	}
	return files.Save(dir, data)
}
