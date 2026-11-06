// Package catalog owns products, courses and prices (schema "catalog").
// Other modules use it only through the exported methods on *Module, e.g. Price.
package catalog

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"easygym/internal/platform"
)

// JSON names match the frontend's card and detail shapes.
type Option struct {
	Label string `json:"label"`
	Price int64  `json:"price"`
}

type QA struct {
	Q string `json:"q"`
	A string `json:"a"`
}

type Category struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Product struct {
	Slug     string   `json:"slug"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Tag      string   `json:"tag"`
	Image    string   `json:"image"`
	Reviews  int      `json:"reviews"`
	Variants []Option `json:"variants"`
}

type ProductDetail struct {
	Product
	Images    []string `json:"images"`
	Desc      string   `json:"desc"`
	Detail    string   `json:"detail"`
	Questions []QA     `json:"questions"`
}

type Course struct {
	Slug  string   `json:"slug"`
	Name  string   `json:"name"`
	Mode  string   `json:"mode"`
	Card  string   `json:"card"`
	Image string   `json:"image"`
	Desc  string   `json:"desc"`
	Plans []Option `json:"plans"`
}

type CourseDetail struct {
	Course
	Body string `json:"body"`
}

type Catalog struct {
	Categories []Category `json:"categories"`
	Products   []Product  `json:"products"`
	Courses    []Course   `json:"courses"`
}

type FAQGroup struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Items []QA   `json:"items"`
}

type Module struct {
	db   *pgxpool.Pool
	auth *platform.Auth
}

func New(db *pgxpool.Pool, auth *platform.Auth) *Module { return &Module{db: db, auth: auth} }

func (m *Module) Register(mux *http.ServeMux) {
	m.registerAdmin(mux)
	mux.HandleFunc("GET /api/catalog", func(w http.ResponseWriter, r *http.Request) {
		c, err := m.Catalog(r.Context())
		platform.Respond(w, r, c, err)
	})
	mux.HandleFunc("GET /api/products/{slug}", func(w http.ResponseWriter, r *http.Request) {
		p, err := m.Product(r.Context(), r.PathValue("slug"))
		platform.Respond(w, r, p, err)
	})
	mux.HandleFunc("GET /api/courses/{slug}", func(w http.ResponseWriter, r *http.Request) {
		c, err := m.Course(r.Context(), r.PathValue("slug"))
		platform.Respond(w, r, c, err)
	})
	mux.HandleFunc("GET /api/faq", func(w http.ResponseWriter, r *http.Request) {
		f, err := m.FAQ(r.Context())
		platform.Respond(w, r, f, err)
	})
}

const productCols = `p.slug, p.name, p.category, p.tag, p.image, p.reviews,
	coalesce((select json_agg(json_build_object('label', v.label, 'price', v.price) order by v.position)
		from catalog.product_variants v where v.product_slug = p.slug), '[]')`

const courseCols = `c.slug, c.name, c.mode, c.card, c.image, c.summary,
	coalesce((select json_agg(json_build_object('label', x.label, 'price', x.price) order by x.position)
		from catalog.course_plans x where x.course_slug = c.slug), '[]')`

func scanProduct(row pgx.Row, extra ...any) (Product, error) {
	var p Product
	err := row.Scan(append([]any{&p.Slug, &p.Name, &p.Category, &p.Tag, &p.Image, &p.Reviews, &p.Variants}, extra...)...)
	return p, err
}

func scanCourse(row pgx.Row, extra ...any) (Course, error) {
	var c Course
	err := row.Scan(append([]any{&c.Slug, &c.Name, &c.Mode, &c.Card, &c.Image, &c.Desc, &c.Plans}, extra...)...)
	return c, err
}

// Catalog returns every active card: small enough (tens of rows) to send in one response.
func (m *Module) Catalog(ctx context.Context) (Catalog, error) {
	var out Catalog
	rows, _ := m.db.Query(ctx, `select slug, label from catalog.categories order by position`)
	cats, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Category, error) {
		var c Category
		return c, r.Scan(&c.Value, &c.Label)
	})
	if err != nil {
		return out, fmt.Errorf("categories: %w", err)
	}
	rows, _ = m.db.Query(ctx, `select `+productCols+` from catalog.products p where p.active order by p.position`)
	products, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Product, error) { return scanProduct(r) })
	if err != nil {
		return out, fmt.Errorf("products: %w", err)
	}
	rows, _ = m.db.Query(ctx, `select `+courseCols+` from catalog.courses c where c.active order by c.position`)
	courses, err := pgx.CollectRows(rows, func(r pgx.CollectableRow) (Course, error) { return scanCourse(r) })
	if err != nil {
		return out, fmt.Errorf("courses: %w", err)
	}
	return Catalog{Categories: cats, Products: products, Courses: courses}, nil
}

func (m *Module) Product(ctx context.Context, slug string) (ProductDetail, error) {
	var d ProductDetail
	p, err := scanProduct(m.db.QueryRow(ctx, `select `+productCols+`, p.images, p.description, p.detail,
		coalesce((select json_agg(json_build_object('q', q.question, 'a', q.answer) order by q.position)
			from catalog.product_questions q where q.product_slug = p.slug), '[]')
		from catalog.products p where p.slug = $1 and p.active`, slug), &d.Images, &d.Desc, &d.Detail, &d.Questions)
	if errors.Is(err, pgx.ErrNoRows) {
		return d, platform.ErrNotFound
	}
	d.Product = p
	return d, err
}

func (m *Module) Course(ctx context.Context, slug string) (CourseDetail, error) {
	var d CourseDetail
	c, err := scanCourse(m.db.QueryRow(ctx, `select `+courseCols+`, c.body from catalog.courses c where c.slug = $1 and c.active`, slug), &d.Body)
	if errors.Is(err, pgx.ErrNoRows) {
		return d, platform.ErrNotFound
	}
	d.Course = c
	return d, err
}

// FAQ groups each product's questions, in catalog order.
func (m *Module) FAQ(ctx context.Context) ([]FAQGroup, error) {
	rows, _ := m.db.Query(ctx, `select p.slug, p.name,
		json_agg(json_build_object('q', q.question, 'a', q.answer) order by q.position)
		from catalog.products p join catalog.product_questions q on q.product_slug = p.slug
		where p.active group by p.slug, p.name, p.position order by p.position`)
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (FAQGroup, error) {
		var g FAQGroup
		return g, r.Scan(&g.Slug, &g.Name, &g.Items)
	})
}

// Price is the catalog's public API for the orders module: the current name and price of one variant.
// It returns platform.ErrNotFound when the product is inactive or the variant does not exist.
func (m *Module) Price(ctx context.Context, slug, variant string) (name string, price int64, err error) {
	err = m.db.QueryRow(ctx, `select p.name, v.price from catalog.products p
		join catalog.product_variants v on v.product_slug = p.slug
		where p.slug = $1 and v.label = $2 and p.active`, slug, variant).Scan(&name, &price)
	if errors.Is(err, pgx.ErrNoRows) {
		err = platform.ErrNotFound
	}
	return name, price, err
}
