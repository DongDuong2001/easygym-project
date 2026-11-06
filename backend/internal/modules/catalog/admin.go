package catalog

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"easygym/internal/platform"
)

// ProductInput is a whole product as the admin page edits it. Hiding a product is active=false:
// rows are never deleted, so old orders keep making sense.
type ProductInput struct {
	Slug      string   `json:"slug"` // filled from the URL on save
	Name      string   `json:"name"`
	Category  string   `json:"category"`
	Tag       string   `json:"tag"`
	Image     string   `json:"image"`
	Images    []string `json:"images"`
	Reviews   int      `json:"reviews"`
	Desc      string   `json:"desc"`
	Detail    string   `json:"detail"` // HTML, written by staff
	Position  int      `json:"position"`
	Active    bool     `json:"active"`
	Variants  []Option `json:"variants"`
	Questions []QA     `json:"questions"`
}

type CourseInput struct {
	Slug     string   `json:"slug"`
	Name     string   `json:"name"`
	Mode     string   `json:"mode"`
	Card     string   `json:"card"` // card photo
	Image    string   `json:"image"`
	Desc     string   `json:"desc"`
	Body     string   `json:"body"` // HTML, written by staff
	Position int      `json:"position"`
	Active   bool     `json:"active"`
	Plans    []Option `json:"plans"`
}

func (m *Module) registerAdmin(mux *http.ServeMux) {
	staff := []string{"staff", "admin"}
	mux.HandleFunc("GET /api/admin/products", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		rows, _ := m.db.Query(r.Context(), `select p.slug, p.name, p.category, p.tag, p.image, p.images, p.reviews,
			p.description, p.detail, p.position, p.active,
			coalesce((select json_agg(json_build_object('label', v.label, 'price', v.price) order by v.position)
				from catalog.product_variants v where v.product_slug = p.slug), '[]'),
			coalesce((select json_agg(json_build_object('q', q.question, 'a', q.answer) order by q.position)
				from catalog.product_questions q where q.product_slug = p.slug), '[]')
			from catalog.products p order by p.position, p.slug`)
		out, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (ProductInput, error) {
			var p ProductInput
			return p, row.Scan(&p.Slug, &p.Name, &p.Category, &p.Tag, &p.Image, &p.Images, &p.Reviews,
				&p.Desc, &p.Detail, &p.Position, &p.Active, &p.Variants, &p.Questions)
		})
		platform.Respond(w, r, out, err)
	}, staff...))
	mux.HandleFunc("PUT /api/admin/products/{slug}", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		var in ProductInput
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		in.Slug = r.PathValue("slug")
		platform.Respond(w, r, map[string]string{"slug": in.Slug}, m.SaveProduct(r.Context(), in))
	}, staff...))
	mux.HandleFunc("GET /api/admin/courses", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		rows, _ := m.db.Query(r.Context(), `select c.slug, c.name, c.mode, c.card, c.image, c.summary, c.body, c.position, c.active,
			coalesce((select json_agg(json_build_object('label', x.label, 'price', x.price) order by x.position)
				from catalog.course_plans x where x.course_slug = c.slug), '[]')
			from catalog.courses c order by c.position, c.slug`)
		out, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (CourseInput, error) {
			var c CourseInput
			return c, row.Scan(&c.Slug, &c.Name, &c.Mode, &c.Card, &c.Image, &c.Desc, &c.Body, &c.Position, &c.Active, &c.Plans)
		})
		platform.Respond(w, r, out, err)
	}, staff...))
	mux.HandleFunc("PUT /api/admin/courses/{slug}", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		var in CourseInput
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		in.Slug = r.PathValue("slug")
		platform.Respond(w, r, map[string]string{"slug": in.Slug}, m.SaveCourse(r.Context(), in))
	}, staff...))
	mux.HandleFunc("PUT /api/admin/categories/{slug}", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		var in struct {
			Label    string `json:"label"`
			Position int    `json:"position"`
		}
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		slug := r.PathValue("slug")
		label, err := platform.Text("label", "tên danh mục", in.Label, true, 60)
		if err == nil {
			err = platform.Slug("slug", slug)
		}
		if err == nil {
			_, err = m.db.Exec(r.Context(), `insert into catalog.categories (slug, label, position) values ($1, $2, $3)
				on conflict (slug) do update set label = excluded.label, position = excluded.position`, slug, label, in.Position)
		}
		platform.Respond(w, r, map[string]string{"slug": slug}, err)
	}, staff...))
}

// CheckProduct trims and validates a product before it is saved.
func CheckProduct(in ProductInput) (ProductInput, error) {
	var err error
	if err = platform.Slug("slug", in.Slug); err != nil {
		return in, err
	}
	if err = platform.Slug("category", in.Category); err != nil {
		return in, err
	}
	if in.Name, err = platform.Text("name", "tên sản phẩm", in.Name, true, 200); err != nil {
		return in, err
	}
	if in.Tag, err = platform.Text("tag", "nhãn", in.Tag, false, 40); err != nil {
		return in, err
	}
	if in.Image, err = platform.ImageURL("image", in.Image, true); err != nil {
		return in, err
	}
	if len(in.Images) > 30 {
		return in, platform.ValidationError{Field: "images", Message: "Tối đa 30 ảnh"}
	}
	for i := range in.Images {
		if in.Images[i], err = platform.ImageURL("images", in.Images[i], true); err != nil {
			return in, err
		}
	}
	if in.Reviews < 0 || in.Reviews > 1_000_000 {
		return in, platform.ValidationError{Field: "reviews", Message: "Số đánh giá không hợp lệ"}
	}
	if in.Desc, err = platform.Text("desc", "mô tả", in.Desc, false, 3000); err != nil {
		return in, err
	}
	if in.Detail, err = platform.Text("detail", "chi tiết", in.Detail, false, 50_000); err != nil {
		return in, err
	}
	if in.Variants, err = checkOptions("variants", in.Variants); err != nil {
		return in, err
	}
	if len(in.Questions) > 50 {
		return in, platform.ValidationError{Field: "questions", Message: "Tối đa 50 câu hỏi"}
	}
	for i, q := range in.Questions {
		if in.Questions[i].Q, err = platform.Text("questions", "câu hỏi", q.Q, true, 300); err != nil {
			return in, err
		}
		if in.Questions[i].A, err = platform.Text("questions", "câu trả lời", q.A, true, 3000); err != nil {
			return in, err
		}
	}
	return in, nil
}

func CheckCourse(in CourseInput) (CourseInput, error) {
	var err error
	if err = platform.Slug("slug", in.Slug); err != nil {
		return in, err
	}
	if in.Name, err = platform.Text("name", "tên khóa tập", in.Name, true, 200); err != nil {
		return in, err
	}
	if in.Mode, err = platform.Text("mode", "hình thức", in.Mode, false, 60); err != nil {
		return in, err
	}
	if in.Card, err = platform.ImageURL("card", in.Card, true); err != nil {
		return in, err
	}
	if in.Image, err = platform.ImageURL("image", in.Image, true); err != nil {
		return in, err
	}
	if in.Desc, err = platform.Text("desc", "mô tả", in.Desc, false, 3000); err != nil {
		return in, err
	}
	if in.Body, err = platform.Text("body", "nội dung", in.Body, false, 50_000); err != nil {
		return in, err
	}
	in.Plans, err = checkOptions("plans", in.Plans)
	return in, err
}

// checkOptions validates the price list: 1 to 30 options, each label unique.
func checkOptions(field string, opts []Option) ([]Option, error) {
	if len(opts) == 0 || len(opts) > 30 {
		return opts, platform.ValidationError{Field: field, Message: "Cần từ 1 đến 30 lựa chọn giá"}
	}
	seen := map[string]bool{}
	for i, o := range opts {
		label, err := platform.Text(field, "tên lựa chọn", o.Label, true, 60)
		if err != nil {
			return opts, err
		}
		if seen[label] {
			return opts, platform.ValidationError{Field: field, Message: fmt.Sprintf("Lựa chọn %q bị trùng", label)}
		}
		if o.Price < 0 || o.Price > 10_000_000_000 {
			return opts, platform.ValidationError{Field: field, Message: "Giá không hợp lệ"}
		}
		seen[label] = true
		opts[i].Label = label
	}
	return opts, nil
}

// SaveProduct creates or replaces a product with its variants and questions in one transaction.
func (m *Module) SaveProduct(ctx context.Context, in ProductInput) error {
	in, err := CheckProduct(in)
	if err != nil {
		return err
	}
	labels, prices := splitOptions(in.Variants)
	qs, as := make([]string, len(in.Questions)), make([]string, len(in.Questions))
	for i, q := range in.Questions {
		qs[i], as[i] = q.Q, q.A
	}
	if in.Images == nil {
		in.Images = []string{}
	}
	err = platform.InTx(ctx, m.db, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `insert into catalog.products
			(slug, name, category, tag, image, images, reviews, description, detail, position, active)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			on conflict (slug) do update set name = excluded.name, category = excluded.category, tag = excluded.tag,
				image = excluded.image, images = excluded.images, reviews = excluded.reviews,
				description = excluded.description, detail = excluded.detail, position = excluded.position,
				active = excluded.active, updated_at = now()`,
			in.Slug, in.Name, in.Category, in.Tag, in.Image, in.Images, in.Reviews, in.Desc, in.Detail, in.Position, in.Active); err != nil {
			return err
		}
		batch := &pgx.Batch{}
		batch.Queue(`delete from catalog.product_variants where product_slug = $1`, in.Slug)
		batch.Queue(`insert into catalog.product_variants (product_slug, label, price, position)
			select $1, l, p, i from unnest($2::text[], $3::bigint[]) with ordinality as t(l, p, i)`, in.Slug, labels, prices)
		batch.Queue(`delete from catalog.product_questions where product_slug = $1`, in.Slug)
		batch.Queue(`insert into catalog.product_questions (product_slug, position, question, answer)
			select $1, i, q, a from unnest($2::text[], $3::text[]) with ordinality as t(q, a, i)`, in.Slug, qs, as)
		return tx.SendBatch(ctx, batch).Close()
	})
	return categoryError(err)
}

func (m *Module) SaveCourse(ctx context.Context, in CourseInput) error {
	in, err := CheckCourse(in)
	if err != nil {
		return err
	}
	labels, prices := splitOptions(in.Plans)
	return platform.InTx(ctx, m.db, func(tx pgx.Tx) error {
		batch := &pgx.Batch{}
		batch.Queue(`insert into catalog.courses (slug, name, mode, card, image, summary, body, position, active)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			on conflict (slug) do update set name = excluded.name, mode = excluded.mode, card = excluded.card,
				image = excluded.image, summary = excluded.summary, body = excluded.body,
				position = excluded.position, active = excluded.active, updated_at = now()`,
			in.Slug, in.Name, in.Mode, in.Card, in.Image, in.Desc, in.Body, in.Position, in.Active)
		batch.Queue(`delete from catalog.course_plans where course_slug = $1`, in.Slug)
		batch.Queue(`insert into catalog.course_plans (course_slug, label, price, position)
			select $1, l, p, i from unnest($2::text[], $3::bigint[]) with ordinality as t(l, p, i)`, in.Slug, labels, prices)
		return tx.SendBatch(ctx, batch).Close()
	})
}

func splitOptions(opts []Option) ([]string, []int64) {
	labels, prices := make([]string, len(opts)), make([]int64, len(opts))
	for i, o := range opts {
		labels[i], prices[i] = o.Label, o.Price
	}
	return labels, prices
}

// categoryError turns "category does not exist" (a foreign key violation) into a message for the form.
func categoryError(err error) error {
	var pg *pgconn.PgError
	if errors.As(err, &pg) && pg.Code == "23503" {
		return platform.ValidationError{Field: "category", Message: "Danh mục không tồn tại"}
	}
	return err
}
