package content

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	"easygym/internal/platform"
)

// PostInput is a post as the admin page edits it. Unpublishing (published=false) hides it from the site.
type PostInput struct {
	Slug      string `json:"slug"` // filled from the URL on save
	Title     string `json:"title"`
	Excerpt   string `json:"excerpt"`
	Image     string `json:"image"`
	Cat       string `json:"cat"`
	Date      string `json:"date"` // YYYY-MM-DD
	Body      string `json:"body"` // HTML, written by staff
	Published bool   `json:"published"`
}

func (m *Module) registerAdmin(mux *http.ServeMux) {
	staff := []string{"staff", "admin"}
	mux.HandleFunc("GET /api/admin/posts", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		posts, err := m.adminPosts(r.Context(), "")
		platform.Respond(w, r, posts, err)
	}, staff...))
	mux.HandleFunc("GET /api/admin/posts/{slug}", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		posts, err := m.adminPosts(r.Context(), r.PathValue("slug"))
		if err == nil && len(posts) == 0 {
			err = platform.ErrNotFound
		}
		if err != nil {
			platform.WriteError(w, r, err)
			return
		}
		platform.WriteJSON(w, http.StatusOK, posts[0])
	}, staff...))
	mux.HandleFunc("PUT /api/admin/posts/{slug}", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		var in PostInput
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		in.Slug = r.PathValue("slug")
		platform.Respond(w, r, map[string]string{"slug": in.Slug}, m.SavePost(r.Context(), in))
	}, staff...))
	mux.HandleFunc("PUT /api/admin/pages/{slug}", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		var in struct {
			Body string `json:"body"`
		}
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		slug := r.PathValue("slug")
		err := platform.Slug("slug", slug)
		if err == nil {
			in.Body, err = platform.Text("body", "nội dung", in.Body, true, 50_000)
		}
		if err == nil {
			_, err = m.db.Exec(r.Context(), `insert into content.pages (slug, body) values ($1, $2)
				on conflict (slug) do update set body = excluded.body, updated_at = now()`, slug, in.Body)
		}
		platform.Respond(w, r, map[string]string{"slug": slug}, err)
	}, staff...))
}

// adminPosts lists every post, published or not, newest first; with a slug, the body of that one post.
func (m *Module) adminPosts(ctx context.Context, slug string) ([]PostInput, error) {
	rows, _ := m.db.Query(ctx, `select slug, title, excerpt, image, category, to_char(published_on, 'YYYY-MM-DD'),
		case when $1 = '' then '' else body end, published
		from content.posts where $1 = '' or slug = $1 order by published_on desc, slug`, slug)
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (PostInput, error) {
		var p PostInput
		return p, r.Scan(&p.Slug, &p.Title, &p.Excerpt, &p.Image, &p.Cat, &p.Date, &p.Body, &p.Published)
	})
}

// CheckPost trims and validates a post before it is saved.
func CheckPost(in PostInput) (PostInput, error) {
	var err error
	if err = platform.Slug("slug", in.Slug); err != nil {
		return in, err
	}
	if err = platform.Slug("cat", in.Cat); err != nil {
		return in, err
	}
	if in.Title, err = platform.Text("title", "tiêu đề", in.Title, true, 300); err != nil {
		return in, err
	}
	if in.Excerpt, err = platform.Text("excerpt", "tóm tắt", in.Excerpt, false, 1000); err != nil {
		return in, err
	}
	if in.Image, err = platform.ImageURL("image", in.Image, false); err != nil {
		return in, err
	}
	if _, err := time.Parse(time.DateOnly, in.Date); err != nil {
		return in, platform.ValidationError{Field: "date", Message: "Ngày đăng có dạng năm-tháng-ngày"}
	}
	if in.Body, err = platform.Text("body", "nội dung", in.Body, false, 50_000); err != nil {
		return in, err
	}
	return in, nil
}

func (m *Module) SavePost(ctx context.Context, in PostInput) error {
	in, err := CheckPost(in)
	if err != nil {
		return err
	}
	_, err = m.db.Exec(ctx, `insert into content.posts (slug, title, excerpt, image, category, published_on, body, published)
		values ($1, $2, $3, $4, $5, $6::date, $7, $8)
		on conflict (slug) do update set title = excluded.title, excerpt = excluded.excerpt, image = excluded.image,
			category = excluded.category, published_on = excluded.published_on, body = excluded.body,
			published = excluded.published, updated_at = now()`,
		in.Slug, in.Title, in.Excerpt, in.Image, in.Cat, in.Date, in.Body, in.Published)
	return err
}
