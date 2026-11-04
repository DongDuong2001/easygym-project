package orders

import (
	"context"
	"errors"
	"testing"

	"easygym/internal/platform"
)

type fakePrices map[string]int64

func (f fakePrices) Price(_ context.Context, slug, variant string) (string, int64, error) {
	p, ok := f[slug+"/"+variant]
	if !ok {
		return "", 0, platform.ErrNotFound
	}
	return "Tên " + slug, p, nil
}

func TestBuild(t *testing.T) {
	prices := fakePrices{"kefir/500ml": 85000, "ao/M": 429000}
	ok := Input{Name: " An ", Phone: "0905.152.949", Address: "Gò Vấp", Pay: "cod",
		Lines: []LineInput{{Slug: "kefir", Label: "500ml", Qty: 2}, {Slug: "ao", Label: "M", Qty: 1}}}

	o, err := Build(context.Background(), prices, ok)
	if err != nil {
		t.Fatal(err)
	}
	if o.Total != 2*85000+429000 || o.Name != "An" || o.Phone != "0905152949" || o.Lines[1].Name != "Tên ao" {
		t.Fatalf("got %+v", o)
	}

	bad := map[string]func(*Input){
		"phone":     func(in *Input) { in.Phone = "12345" },
		"address":   func(in *Input) { in.Address = "  " },
		"pay":       func(in *Input) { in.Pay = "card" },
		"empty":     func(in *Input) { in.Lines = nil },
		"qty":       func(in *Input) { in.Lines = []LineInput{{Slug: "kefir", Label: "500ml", Qty: 100}} },
		"gone":      func(in *Input) { in.Lines = []LineInput{{Slug: "kefir", Label: "2L", Qty: 1}} },
		"long name": func(in *Input) { in.Name = string(make([]rune, 121)) },
	}
	for name, mutate := range bad {
		in := ok
		mutate(&in)
		var ve platform.ValidationError
		if _, err := Build(context.Background(), prices, in); !errors.As(err, &ve) {
			t.Errorf("%s: want validation error, got %v", name, err)
		}
	}
}
