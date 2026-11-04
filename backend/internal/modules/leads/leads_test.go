package leads

import "testing"

func TestValidate(t *testing.T) {
	in, err := Validate(Input{Kind: "register", Name: " Bình ", Phone: "0905 152 949", Email: "x@y.z", Course: "group-training"})
	if err != nil || in.Name != "Bình" || in.Phone != "0905152949" || in.Email != "" {
		t.Fatalf("register: %+v %v", in, err)
	}
	in, err = Validate(Input{Kind: "newsletter", Email: " a@b.vn ", Name: "dropped"})
	if err != nil || in.Email != "a@b.vn" || in.Name != "" {
		t.Fatalf("newsletter: %+v %v", in, err)
	}
	for _, bad := range []Input{
		{Kind: "register", Name: "", Phone: "0905152949"},
		{Kind: "register", Name: "A", Phone: "905152949"},
		{Kind: "newsletter", Email: "not-an-email"},
		{Kind: "spam"},
	} {
		if _, err := Validate(bad); err == nil {
			t.Errorf("want error for %+v", bad)
		}
	}
}
