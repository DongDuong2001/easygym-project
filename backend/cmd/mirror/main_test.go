package main

import (
	"slices"
	"testing"
)

func TestFindURLs(t *testing.T) {
	html := `<p><img src="https://easygym.vn/uploads/source/food/a-b_1.webp" alt="x">
		<a href="https://easygym.vn/khoa-tap">link</a> url(https://www.easygym.vn/uploads/medium/s/B.JPG)
		{https://easygym.vn/uploads/source/x.png,"https://easygym.vn/uploads/source/y.jpeg"}
		{https://easygym.vn/uploads/medium/p1.webp,https://easygym.vn/uploads/medium/p2.webp}
		https://other.vn/a.png https://easygym.vn/uploads/slide/intro.mp4 https://easygym.vn.evil.com/a.png</p>`
	want := []string{
		"https://easygym.vn/uploads/source/food/a-b_1.webp",
		"https://www.easygym.vn/uploads/medium/s/B.JPG",
		"https://easygym.vn/uploads/source/x.png",
		"https://easygym.vn/uploads/source/y.jpeg",
		"https://easygym.vn/uploads/medium/p1.webp",
		"https://easygym.vn/uploads/medium/p2.webp",
	}
	if got := FindURLs(html); !slices.Equal(got, want) {
		t.Fatalf("got  %q\nwant %q", got, want)
	}
}
