// Command seed loads the content copied from easygym.vn (db/seed/seed.json) into an empty database.
// It refuses to run on a database that already has products unless -force is given, because -force
// replaces the catalog and content tables and would wipe edits made since.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"

	"easygym/db"
	"easygym/internal/platform"
)

type option struct {
	Label string `json:"label"`
	Price int64  `json:"price"`
}

type seedData struct {
	Categories []struct{ Slug, Label string } `json:"categories"`
	Products   []struct {
		Slug, Name, Category, Image, Tag, Desc, Detail string
		Reviews                                        int
		Variants                                       []option
		Images                                         []string
		Questions                                      []struct{ Q, A string }
	} `json:"products"`
	Courses []struct {
		Slug, Name, Mode, Card, Image, Desc, Body string
		Plans                                     []option
	} `json:"courses"`
	Posts []struct {
		Slug, Title, Excerpt, Image, Date, Cat, Body string
	} `json:"posts"`
	Albums []struct {
		ID, Title string
		Images    []string
	} `json:"albums"`
	Videos []struct{ ID, Title string }  `json:"videos"`
	Pages  []struct{ Slug, Body string } `json:"pages"`
}

func main() {
	force := flag.Bool("force", false, "replace existing catalog and content")
	flag.Parse()
	if err := run(*force); err != nil {
		fmt.Fprintln(os.Stderr, "seed:", err)
		os.Exit(1)
	}
}

func run(force bool) error {
	var s seedData
	if err := json.Unmarshal(db.Seed, &s); err != nil {
		return err
	}
	cfg, err := platform.LoadConfig()
	if err != nil {
		return err
	}
	ctx := context.Background()
	pool, err := platform.OpenDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	var n int
	if err := pool.QueryRow(ctx, `select count(*) from catalog.products`).Scan(&n); err != nil {
		return fmt.Errorf("%w (run migrate first?)", err)
	}
	if n > 0 && !force {
		return errors.New("database already has products; use -force to replace catalog and content")
	}

	return platform.InTx(ctx, pool, func(tx pgx.Tx) error {
		b := &pgx.Batch{}
		q := func(sql string, args ...any) { b.Queue(sql, args...) }
		q(`truncate catalog.categories, catalog.products, catalog.product_variants, catalog.product_questions,
			catalog.courses, catalog.course_plans, content.posts, content.albums, content.album_images,
			content.videos, content.pages`)
		for i, c := range s.Categories {
			q(`insert into catalog.categories (slug, label, position) values ($1, $2, $3)`, c.Slug, c.Label, i)
		}
		for i, p := range s.Products {
			q(`insert into catalog.products (slug, name, category, tag, image, images, reviews, description, detail, position)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
				p.Slug, p.Name, p.Category, p.Tag, p.Image, p.Images, p.Reviews, p.Desc, p.Detail, i)
			for j, v := range p.Variants {
				q(`insert into catalog.product_variants (product_slug, label, price, position) values ($1, $2, $3, $4)`, p.Slug, v.Label, v.Price, j)
			}
			for j, qa := range p.Questions {
				q(`insert into catalog.product_questions (product_slug, position, question, answer) values ($1, $2, $3, $4)`, p.Slug, j, qa.Q, qa.A)
			}
		}
		for i, c := range s.Courses {
			q(`insert into catalog.courses (slug, name, mode, card, image, summary, body, position)
				values ($1, $2, $3, $4, $5, $6, $7, $8)`, c.Slug, c.Name, c.Mode, c.Card, c.Image, c.Desc, c.Body, i)
			for j, pl := range c.Plans {
				q(`insert into catalog.course_plans (course_slug, label, price, position) values ($1, $2, $3, $4)`, c.Slug, pl.Label, pl.Price, j)
			}
		}
		for _, p := range s.Posts {
			day, err := time.Parse("02/01/2006", p.Date)
			if err != nil {
				return fmt.Errorf("post %s: %w", p.Slug, err)
			}
			q(`insert into content.posts (slug, title, excerpt, image, category, published_on, body)
				values ($1, $2, $3, $4, $5, $6, $7)`, p.Slug, p.Title, p.Excerpt, p.Image, p.Cat, day, p.Body)
		}
		for i, a := range s.Albums {
			q(`insert into content.albums (slug, title, position) values ($1, $2, $3)`, a.ID, a.Title, i)
			for j, url := range a.Images {
				q(`insert into content.album_images (album_slug, position, url) values ($1, $2, $3)`, a.ID, j, url)
			}
		}
		for i, v := range s.Videos {
			q(`insert into content.videos (youtube_id, title, position) values ($1, $2, $3)`, v.ID, v.Title, i)
		}
		for _, p := range s.Pages {
			q(`insert into content.pages (slug, body) values ($1, $2)`, p.Slug, p.Body)
		}
		if err := tx.SendBatch(ctx, b).Close(); err != nil {
			return err
		}
		fmt.Printf("seeded %d products, %d courses, %d posts, %d albums, %d videos\n",
			len(s.Products), len(s.Courses), len(s.Posts), len(s.Albums), len(s.Videos))
		return nil
	})
}
