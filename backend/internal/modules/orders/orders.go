// Package orders turns a cart into an order (schema "orders"). Prices always come from the
// catalog through Pricer, never from the browser.
package orders

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"easygym/internal/platform"
)

// Pricer is the slice of the catalog module that orders needs; *catalog.Module satisfies it.
type Pricer interface {
	Price(ctx context.Context, slug, variant string) (name string, price int64, err error)
}

type LineInput struct {
	Slug  string `json:"slug"`
	Label string `json:"label"`
	Qty   int    `json:"qty"`
}

type Input struct {
	Name    string      `json:"name"`
	Phone   string      `json:"phone"`
	Address string      `json:"address"`
	Pay     string      `json:"pay"` // "cod" or "bank"
	Note    string      `json:"note"`
	Lines   []LineInput `json:"lines"`
}

type Line struct {
	Slug, Name, Label string
	Price             int64
	Qty               int
}

type Order struct {
	Input
	Lines []Line
	Total int64
}

// Statuses an order moves through; the admin page changes them.
var Statuses = []string{"new", "confirmed", "shipping", "done", "cancelled"}

// View is a saved order as the customer's history and the admin page show it.
type View struct {
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Address   string     `json:"address"`
	Pay       string     `json:"pay"`
	Note      string     `json:"note"`
	Total     int64      `json:"total"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	Lines     []ViewLine `json:"lines"`
}

type ViewLine struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Price int64  `json:"price"`
	Qty   int    `json:"qty"`
}

type Module struct {
	db     *pgxpool.Pool
	prices Pricer
	auth   *platform.Auth
}

func New(db *pgxpool.Pool, prices Pricer, auth *platform.Auth) *Module {
	return &Module{db: db, prices: prices, auth: auth}
}

func (m *Module) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/orders", func(w http.ResponseWriter, r *http.Request) {
		var in Input
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		// Signing in is optional for checkout: a bad or missing token just places a guest order.
		u, _ := m.auth.User(r)
		code, total, err := m.Place(r.Context(), in, u.ID)
		if err != nil {
			platform.WriteError(w, r, err)
			return
		}
		platform.WriteJSON(w, http.StatusCreated, map[string]any{"code": code, "total": total})
	})
	mux.HandleFunc("GET /api/me/orders", m.auth.Require(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		v, err := m.list(r.Context(), `o.user_id = $1`, u.ID)
		platform.Respond(w, r, v, err)
	}))
	staff := []string{"staff", "admin"}
	mux.HandleFunc("GET /api/admin/orders", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		v, err := m.list(r.Context(), `$1 = '' or o.status = $1`, r.URL.Query().Get("status"))
		platform.Respond(w, r, v, err)
	}, staff...))
	mux.HandleFunc("GET /api/admin/orders/{code}", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		v, err := m.list(r.Context(), `o.code = $1`, r.PathValue("code"))
		if err == nil && len(v) == 0 {
			err = platform.ErrNotFound
		}
		if err != nil {
			platform.WriteError(w, r, err)
			return
		}
		platform.WriteJSON(w, http.StatusOK, v[0])
	}, staff...))
	mux.HandleFunc("PATCH /api/admin/orders/{code}", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		var in struct {
			Status string `json:"status"`
		}
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		err := platform.SetStatus(r.Context(), m.db, "orders.orders", r.PathValue("code"), in.Status, Statuses)
		platform.Respond(w, r, map[string]string{"status": in.Status}, err)
	}, staff...))
}

// list returns orders matching where (one parameter, $1), newest first, with their lines.
// shortcut: the newest 200 only, no paging; add a cursor when the inbox outgrows that.
func (m *Module) list(ctx context.Context, where string, arg string) ([]View, error) {
	rows, _ := m.db.Query(ctx, `select o.code, o.name, o.phone, o.address, o.payment, o.note, o.total, o.status, o.created_at,
		coalesce((select json_agg(json_build_object('slug', l.product_slug, 'name', l.product_name, 'label', l.variant_label,
			'price', l.unit_price, 'qty', l.qty) order by l.position) from orders.order_lines l where l.order_id = o.id), '[]')
		from orders.orders o where `+where+` order by o.created_at desc limit 200`, arg)
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (View, error) {
		var v View
		return v, r.Scan(&v.Code, &v.Name, &v.Phone, &v.Address, &v.Pay, &v.Note, &v.Total, &v.Status, &v.CreatedAt, &v.Lines)
	})
}

// Place validates and prices the order, then saves it with its lines in one transaction.
// userID is empty for a guest checkout.
func (m *Module) Place(ctx context.Context, in Input, userID string) (code string, total int64, err error) {
	o, err := Build(ctx, m.prices, in)
	if err != nil {
		return "", 0, err
	}
	code = platform.NewCode("EG")
	err = platform.InTx(ctx, m.db, func(tx pgx.Tx) error {
		var id int64
		if err := tx.QueryRow(ctx, `insert into orders.orders (code, name, phone, address, payment, note, total, user_id)
			values ($1, $2, $3, $4, $5, $6, $7, nullif($8, '')) returning id`,
			code, o.Name, o.Phone, o.Address, o.Pay, o.Note, o.Total, userID).Scan(&id); err != nil {
			return err
		}
		rows := make([][]any, len(o.Lines))
		for i, l := range o.Lines {
			rows[i] = []any{id, i, l.Slug, l.Name, l.Label, l.Price, l.Qty}
		}
		_, err := tx.CopyFrom(ctx, pgx.Identifier{"orders", "order_lines"},
			[]string{"order_id", "position", "product_slug", "product_name", "variant_label", "unit_price", "qty"},
			pgx.CopyFromRows(rows))
		return err
	})
	return code, o.Total, err
}

// Build checks the customer fields and prices every line. It has no database access of its own,
// so it is tested with a fake Pricer.
func Build(ctx context.Context, prices Pricer, in Input) (Order, error) {
	o := Order{Input: in}
	var err error
	if o.Name, err = platform.Text("name", "họ và tên", in.Name, true, 120); err != nil {
		return o, err
	}
	phone, ok := platform.NormalizePhone(in.Phone)
	if !ok {
		return o, platform.ValidationError{Field: "phone", Message: "Số điện thoại gồm 10 số, bắt đầu bằng 0"}
	}
	o.Phone = phone
	if o.Address, err = platform.Text("address", "địa chỉ giao hàng", in.Address, true, 300); err != nil {
		return o, err
	}
	if o.Note, err = platform.Text("note", "ghi chú", in.Note, false, 1000); err != nil {
		return o, err
	}
	if in.Pay != "cod" && in.Pay != "bank" {
		return o, platform.ValidationError{Field: "pay", Message: "Chọn cách thanh toán"}
	}
	if len(in.Lines) == 0 || len(in.Lines) > 50 {
		return o, platform.ValidationError{Field: "lines", Message: "Giỏ hàng trống hoặc quá nhiều sản phẩm"}
	}
	for _, l := range in.Lines {
		if l.Qty < 1 || l.Qty > 99 {
			return o, platform.ValidationError{Field: "lines", Message: "Số lượng mỗi sản phẩm từ 1 đến 99"}
		}
		name, price, err := prices.Price(ctx, l.Slug, l.Label)
		if errors.Is(err, platform.ErrNotFound) {
			return o, platform.ValidationError{Field: "lines", Message: "Có sản phẩm trong giỏ không còn bán, vui lòng tải lại giỏ hàng"}
		}
		if err != nil {
			return o, err
		}
		o.Lines = append(o.Lines, Line{Slug: l.Slug, Name: name, Label: l.Label, Price: price, Qty: l.Qty})
		o.Total += price * int64(l.Qty)
	}
	return o, nil
}
