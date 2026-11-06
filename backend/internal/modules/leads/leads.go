// Package leads stores course sign-ups, contact requests and newsletter emails (schema "leads").
package leads

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"easygym/internal/platform"
)

type Input struct {
	Kind   string `json:"kind"` // "register" or "newsletter"
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	Course string `json:"course"`
	Plan   string `json:"plan"`
	Note   string `json:"note"`
	Source string `json:"source"` // page or form it came from
}

// Statuses a lead moves through on the admin page.
var Statuses = []string{"new", "contacted", "done", "spam"}

// Lead is a saved sign-up as the admin page lists it.
type Lead struct {
	Input
	Code      string    `json:"code"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type Module struct {
	db   *pgxpool.Pool
	auth *platform.Auth
}

func New(db *pgxpool.Pool, auth *platform.Auth) *Module { return &Module{db: db, auth: auth} }

func (m *Module) Register(mux *http.ServeMux) {
	staff := []string{"staff", "admin"}
	// shortcut: the newest 500 only, no paging; add a cursor when the inbox outgrows that.
	mux.HandleFunc("GET /api/admin/leads", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		q := r.URL.Query()
		rows, _ := m.db.Query(r.Context(), `select code, kind, name, phone, email, course, plan, note, source, status, created_at
			from leads.leads where ($1 = '' or status = $1) and ($2 = '' or kind = $2)
			order by created_at desc limit 500`, q.Get("status"), q.Get("kind"))
		leads, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Lead, error) {
			var l Lead
			return l, row.Scan(&l.Code, &l.Kind, &l.Name, &l.Phone, &l.Email, &l.Course, &l.Plan, &l.Note, &l.Source, &l.Status, &l.CreatedAt)
		})
		platform.Respond(w, r, leads, err)
	}, staff...))
	mux.HandleFunc("PATCH /api/admin/leads/{code}", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		var in struct {
			Status string `json:"status"`
		}
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		err := platform.SetStatus(r.Context(), m.db, "leads.leads", r.PathValue("code"), in.Status, Statuses)
		platform.Respond(w, r, map[string]string{"status": in.Status}, err)
	}, staff...))

	mux.HandleFunc("POST /api/leads", func(w http.ResponseWriter, r *http.Request) {
		var in Input
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		code, err := m.Create(r.Context(), in)
		if err != nil {
			platform.WriteError(w, r, err)
			return
		}
		platform.WriteJSON(w, http.StatusCreated, map[string]string{"code": code})
	})
}

func (m *Module) Create(ctx context.Context, in Input) (string, error) {
	in, err := Validate(in)
	if err != nil {
		return "", err
	}
	code := platform.NewCode("DK")
	_, err = m.db.Exec(ctx, `insert into leads.leads (code, kind, name, phone, email, course, plan, note, source)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		code, in.Kind, in.Name, in.Phone, in.Email, in.Course, in.Plan, in.Note, in.Source)
	return code, err
}

var emailRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// Validate trims and checks a lead; it returns the cleaned copy or a platform.ValidationError.
func Validate(in Input) (Input, error) {
	var err error
	if in.Source, err = platform.Text("source", "nguồn", in.Source, false, 60); err != nil {
		return in, err
	}
	switch in.Kind {
	case "newsletter":
		if in.Email, err = platform.Text("email", "email", in.Email, true, 254); err != nil {
			return in, err
		}
		if !emailRe.MatchString(in.Email) {
			return in, platform.ValidationError{Field: "email", Message: "Email chưa đúng định dạng"}
		}
		in.Name, in.Phone, in.Course, in.Plan, in.Note = "", "", "", "", ""
		return in, nil
	case "register":
		fields := []struct {
			name, label string
			v           *string
			required    bool
			max         int
		}{
			{"name", "họ và tên", &in.Name, true, 120},
			{"course", "khóa học", &in.Course, false, 120},
			{"plan", "gói", &in.Plan, false, 60},
			{"note", "ghi chú", &in.Note, false, 2000},
		}
		for _, f := range fields {
			if *f.v, err = platform.Text(f.name, f.label, *f.v, f.required, f.max); err != nil {
				return in, err
			}
		}
		phone, ok := platform.NormalizePhone(in.Phone)
		if !ok {
			return in, platform.ValidationError{Field: "phone", Message: "Số điện thoại gồm 10 số, bắt đầu bằng 0"}
		}
		in.Phone, in.Email = phone, ""
		return in, nil
	}
	return in, platform.ValidationError{Field: "kind", Message: "Loại đăng ký không hợp lệ"}
}
