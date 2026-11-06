package accounts

import (
	"context"
	"slices"
	"strings"
	"testing"

	"easygym/internal/platform"
)

func TestCleanWishlist(t *testing.T) {
	got, err := CleanWishlist([]string{"kefir", "ao-tap", "kefir"})
	if err != nil || !slices.Equal(got, []string{"kefir", "ao-tap"}) {
		t.Fatalf("got %v %v", got, err)
	}
	if got, err := CleanWishlist(nil); err != nil || len(got) != 0 {
		t.Fatalf("empty list: %v %v", got, err)
	}
	for _, bad := range [][]string{{"Kefir"}, {"a b"}, {"-x"}, {""}, {strings.Repeat("a", 121)}, make([]string, 201)} {
		if _, err := CleanWishlist(bad); err == nil {
			t.Errorf("%q should fail", bad)
		}
	}
}

func TestSetRoleChecksBeforeTouchingTheDatabase(t *testing.T) {
	m := &Module{} // no db: these must fail before any query
	me := platform.User{ID: "me", Role: "admin"}
	if err := m.SetRole(context.Background(), me, "other", "owner"); err == nil {
		t.Error("unknown role must fail")
	}
	if err := m.SetRole(context.Background(), me, "me", "customer"); err == nil {
		t.Error("an admin must not change their own role")
	}
}
