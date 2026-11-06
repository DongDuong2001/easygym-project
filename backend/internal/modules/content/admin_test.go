package content

import "testing"

func TestCheckPost(t *testing.T) {
	ok := PostInput{Slug: "tap-gi-hom-nay", Title: " Tập gì hôm nay ", Cat: "tin-tuc", Date: "2026-05-30", Body: "<p>Hi</p>"}
	got, err := CheckPost(ok)
	if err != nil || got.Title != "Tập gì hôm nay" {
		t.Fatalf("got %+v %v", got, err)
	}
	bad := map[string]func(*PostInput){
		"date format": func(p *PostInput) { p.Date = "30/05/2026" },
		"bad date":    func(p *PostInput) { p.Date = "2026-02-30" },
		"no title":    func(p *PostInput) { p.Title = "" },
		"category":    func(p *PostInput) { p.Cat = "Tin Tức" },
		"image":       func(p *PostInput) { p.Image = "data:image/png;base64,AAAA" },
	}
	for name, change := range bad {
		in := ok
		change(&in)
		if _, err := CheckPost(in); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}
