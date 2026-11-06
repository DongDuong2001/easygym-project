package catalog

import (
	"strings"
	"testing"
)

func TestCheckProduct(t *testing.T) {
	ok := ProductInput{Slug: "kefir", Name: " Kefir ", Category: "food", Image: "/api/files/abc.webp",
		Images: []string{"https://cdn.example/a.jpg"}, Variants: []Option{{" 500ml ", 85000}, {"1L", 150000}},
		Questions: []QA{{"Bảo quản?", "Ngăn mát"}}}
	got, err := CheckProduct(ok)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Kefir" || got.Variants[0].Label != "500ml" {
		t.Fatalf("not trimmed: %+v", got)
	}

	bad := map[string]func(*ProductInput){
		"slug":         func(p *ProductInput) { p.Slug = "Kefir 1" },
		"category":     func(p *ProductInput) { p.Category = "" },
		"name":         func(p *ProductInput) { p.Name = " " },
		"image scheme": func(p *ProductInput) { p.Image = "http://x/a.jpg" },
		"image js":     func(p *ProductInput) { p.Images = []string{"javascript:alert(1)"} },
		"no variants":  func(p *ProductInput) { p.Variants = nil },
		"dup variant":  func(p *ProductInput) { p.Variants = []Option{{"1L", 1}, {"1L ", 2}} },
		"neg price":    func(p *ProductInput) { p.Variants = []Option{{"1L", -1}} },
		"empty answer": func(p *ProductInput) { p.Questions = []QA{{"Hỏi", " "}} },
		"long detail":  func(p *ProductInput) { p.Detail = strings.Repeat("a", 50_001) },
		"many reviews": func(p *ProductInput) { p.Reviews = -1 },
	}
	for name, change := range bad {
		in := ok
		in.Variants = append([]Option(nil), ok.Variants...)
		change(&in)
		if _, err := CheckProduct(in); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestCheckCourse(t *testing.T) {
	ok := CourseInput{Slug: "tap-luyen-1-1", Name: "Tập 1-1", Card: "https://a/c.jpg", Image: "https://a/i.jpg",
		Plans: []Option{{"1 tháng", 8_700_000}}}
	if _, err := CheckCourse(ok); err != nil {
		t.Fatal(err)
	}
	ok.Plans = nil
	if _, err := CheckCourse(ok); err == nil {
		t.Error("a course needs at least one plan")
	}
}
