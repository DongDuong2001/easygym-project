// Package leads stores course sign-ups, contact requests and newsletter emails (schema "leads").
package leads

import (
	"context"
	"net/http"
	"regexp"

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

type Module struct{ db *pgxpool.Pool }

func New(db *pgxpool.Pool) *Module { return &Module{db: db} }

func (m *Module) Register(mux *http.ServeMux) {
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
