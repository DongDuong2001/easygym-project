// Package accounts owns signed-in users, their role and their wishlist (schema "accounts").
// It implements platform.Roles, so every protected route asks it who is admin or staff.
package accounts

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"easygym/internal/platform"
)

type Account struct {
	platform.User
	CreatedAt time.Time `json:"createdAt"`
	LastSeen  time.Time `json:"lastSeen"`
}

type Module struct {
	db     *pgxpool.Pool
	auth   *platform.Auth
	admins []string // lowercase emails that become admin on sign-in
}

func New(db *pgxpool.Pool, auth *platform.Auth, adminEmails []string) *Module {
	return &Module{db: db, auth: auth, admins: adminEmails}
}

func (m *Module) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/me", m.auth.Require(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		platform.WriteJSON(w, http.StatusOK, u)
	}))
	mux.HandleFunc("GET /api/me/wishlist", m.auth.Require(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		rows, _ := m.db.Query(r.Context(), `select product_slug from accounts.wishlist where user_id = $1 order by position`, u.ID)
		slugs, err := pgx.CollectRows(rows, pgx.RowTo[string])
		platform.Respond(w, r, map[string][]string{"slugs": slugs}, err)
	}))
	mux.HandleFunc("PUT /api/me/wishlist", m.auth.Require(func(w http.ResponseWriter, r *http.Request, u platform.User) {
		var in struct {
			Slugs []string `json:"slugs"`
		}
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		slugs, err := m.SetWishlist(r.Context(), u.ID, in.Slugs)
		platform.Respond(w, r, map[string][]string{"slugs": slugs}, err)
	}))
	mux.HandleFunc("GET /api/admin/users", m.auth.Require(func(w http.ResponseWriter, r *http.Request, _ platform.User) {
		rows, _ := m.db.Query(r.Context(), `select id, email, name, role, created_at, last_seen
			from accounts.users order by role = 'customer', last_seen desc limit 500`)
		users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Account, error) {
			var a Account
			return a, row.Scan(&a.ID, &a.Email, &a.Name, &a.Role, &a.CreatedAt, &a.LastSeen)
		})
		platform.Respond(w, r, users, err)
	}, "admin"))
	mux.HandleFunc("PUT /api/admin/users/{id}/role", m.auth.Require(func(w http.ResponseWriter, r *http.Request, me platform.User) {
		var in struct {
			Role string `json:"role"`
		}
		if err := platform.DecodeJSON(w, r, &in); err != nil {
			platform.WriteError(w, r, err)
			return
		}
		err := m.SetRole(r.Context(), me, r.PathValue("id"), in.Role)
		platform.Respond(w, r, map[string]string{"role": in.Role}, err)
	}, "admin"))
}

// Role records the user on every signed-in request (refreshing email, name and last_seen) and
// returns their role. An email listed in ADMIN_EMAILS is promoted to admin; nobody is demoted here.
func (m *Module) Role(ctx context.Context, u platform.User) (string, error) {
	var role string
	err := m.db.QueryRow(ctx, `insert into accounts.users (id, email, name, role)
		values ($1, $2, $3, case when $4 then 'admin' else 'customer' end)
		on conflict (id) do update set email = excluded.email, name = excluded.name, last_seen = now(),
			role = case when $4 then 'admin' else accounts.users.role end
		returning role`, u.ID, u.Email, u.Name, u.Email != "" && slices.Contains(m.admins, u.Email)).Scan(&role)
	return role, err
}

// SetWishlist replaces the user's wishlist with slugs (deduplicated, in order).
func (m *Module) SetWishlist(ctx context.Context, userID string, slugs []string) ([]string, error) {
	slugs, err := CleanWishlist(slugs)
	if err != nil {
		return nil, err
	}
	return slugs, platform.InTx(ctx, m.db, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `delete from accounts.wishlist where user_id = $1`, userID); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `insert into accounts.wishlist (user_id, product_slug, position)
			select $1, s, i from unnest($2::text[]) with ordinality as t(s, i)`, userID, slugs)
		return err
	})
}

// CleanWishlist checks every slug and drops repeats; at most 200 items.
func CleanWishlist(slugs []string) ([]string, error) {
	if len(slugs) > 200 {
		return nil, platform.ValidationError{Field: "slugs", Message: "Danh sách yêu thích tối đa 200 sản phẩm"}
	}
	out := make([]string, 0, len(slugs))
	for _, s := range slugs {
		if err := platform.Slug("slugs", s); err != nil {
			return nil, err
		}
		if !slices.Contains(out, s) {
			out = append(out, s)
		}
	}
	return out, nil
}

// SetRole changes another user's role. Admins cannot change their own, so the last admin cannot lock everyone out.
func (m *Module) SetRole(ctx context.Context, me platform.User, id, role string) error {
	if !slices.Contains([]string{"customer", "staff", "admin"}, role) {
		return platform.ValidationError{Field: "role", Message: "Quyền không hợp lệ"}
	}
	if id == me.ID {
		return platform.ValidationError{Field: "role", Message: "Không thể tự đổi quyền của chính mình"}
	}
	tag, err := m.db.Exec(ctx, `update accounts.users set role = $2 where id = $1`, id, role)
	if err == nil && tag.RowsAffected() == 0 {
		err = platform.ErrNotFound
	}
	return err
}
