// Package orders turns a cart into an order (schema "orders"). Prices always come from the
// catalog through Pricer, never from the browser.
package orders

import (
	"context"
	"errors"
	"net/http"

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

type Module struct {
	db     *pgxpool.Pool
	prices Pricer
}

func New(db *pgxpool.Pool, prices Pricer) *Module { return &Module{db: db, prices: prices} }

func (m *Module) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/orders", func(w http.ResponseWriter, r *http.Request) {
		var in Input
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		code, total, err := m.Place(r.Context(), in)
		if err != nil {
			platform.WriteError(w, r, err)
			return
		}
		platform.WriteJSON(w, http.StatusCreated, map[string]any{"code": code, "total": total})
	})
}

// Place validates and prices the order, then saves it with its lines in one transaction.
func (m *Module) Place(ctx context.Context, in Input) (code string, total int64, err error) {
	o, err := Build(ctx, m.prices, in)
	if err != nil {
		return "", 0, err
	}
	code = platform.NewCode("EG")
	err = platform.InTx(ctx, m.db, func(tx pgx.Tx) error {
		var id int64
		if err := tx.QueryRow(ctx, `insert into orders.orders (code, name, phone, address, payment, note, total)
			values ($1, $2, $3, $4, $5, $6, $7) returning id`,
			code, o.Name, o.Phone, o.Address, o.Pay, o.Note, o.Total).Scan(&id); err != nil {
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
