// Package content owns news posts, photo albums, videos and static pages (schema "content").
package content

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"easygym/internal/platform"
)

type Post struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Image   string `json:"image"`
	Date    string `json:"date"` // dd/mm/yyyy, as the site prints it
	Cat     string `json:"cat"`
}

type PostDetail struct {
	Post
	Body string `json:"body"`
}

type Album struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Images []string `json:"images"`
}

type Video struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Media is what the home page and the gallery/video pages need in one call.
type Media struct {
	Albums []Album `json:"albums"`
	Videos []Video `json:"videos"`
}

type Module struct {
	db   *pgxpool.Pool
	auth *platform.Auth
}

func New(db *pgxpool.Pool, auth *platform.Auth) *Module { return &Module{db: db, auth: auth} }

func (m *Module) Register(mux *http.ServeMux) {
	m.registerAdmin(mux)
	mux.HandleFunc("GET /api/posts", func(w http.ResponseWriter, r *http.Request) {
		p, err := m.Posts(r.Context(), r.URL.Query().Get("cat"))
		platform.Respond(w, r, p, err)
	})
	mux.HandleFunc("GET /api/posts/{slug}", func(w http.ResponseWriter, r *http.Request) {
		p, err := m.Post(r.Context(), r.PathValue("slug"))
		platform.Respond(w, r, p, err)
	})
	mux.HandleFunc("GET /api/media", func(w http.ResponseWriter, r *http.Request) {
		md, err := m.Media(r.Context())
		platform.Respond(w, r, md, err)
	})
	mux.HandleFunc("GET /api/pages/{slug}", func(w http.ResponseWriter, r *http.Request) {
		var body string
		err := m.db.QueryRow(r.Context(), `select body from content.pages where slug = $1`, r.PathValue("slug")).Scan(&body)
		if errors.Is(err, pgx.ErrNoRows) {
			err = platform.ErrNotFound
		}
		platform.Respond(w, r, map[string]string{"body": body}, err)
	})
}

const postCols = `slug, title, excerpt, image, to_char(published_on, 'DD/MM/YYYY'), category`

// Posts lists published posts, newest first; cat filters by category when set.
func (m *Module) Posts(ctx context.Context, cat string) ([]Post, error) {
	rows, _ := m.db.Query(ctx, `select `+postCols+` from content.posts
		where published and ($1 = '' or category = $1) order by published_on desc, slug`, cat)
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (Post, error) {
		var p Post
		return p, r.Scan(&p.Slug, &p.Title, &p.Excerpt, &p.Image, &p.Date, &p.Cat)
	})
}

func (m *Module) Post(ctx context.Context, slug string) (PostDetail, error) {
	var p PostDetail
	err := m.db.QueryRow(ctx, `select `+postCols+`, body from content.posts where slug = $1 and published`, slug).
		Scan(&p.Slug, &p.Title, &p.Excerpt, &p.Image, &p.Date, &p.Cat, &p.Body)
	if errors.Is(err, pgx.ErrNoRows) {
		err = platform.ErrNotFound
	}
	return p, err
}

func (m *Module) Media(ctx context.Context) (Media, error) {
	var md Media
	rows, _ := m.db.Query(ctx, `select a.slug, a.title,
		coalesce((select array_agg(i.url order by i.position) from content.album_images i where i.album_slug = a.slug), '{}')
		from content.albums a order by a.position`)
	albums, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Album, error) {
		var a Album
		return a, r.Scan(&a.ID, &a.Title, &a.Images)
	})
	if err != nil {
		return md, err
	}
	rows, _ = m.db.Query(ctx, `select youtube_id, title from content.videos order by position`)
	videos, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Video, error) {
		var v Video
		return v, r.Scan(&v.ID, &v.Title)
	})
	return Media{Albums: albums, Videos: videos}, err
}
