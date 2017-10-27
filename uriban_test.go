package uriban

import "testing"
import "net/url"

func TestReplace(t *testing.T) {
	ft := func(s string) string { return "123" }
	cases := []struct {
		url   string
		want  string
		modes []Option
	}{
		{"", "", nil},
		{"url", "url", nil},

		//BEGIN TEST MODES
		{"", "", []Option{WithOption(All, ModeNothing())}},
		{"url", "url", []Option{WithOption(All, ModeNothing())}},
		{"url", "", []Option{WithOption(All, ModeHidden())}},
		{"url", "*****", []Option{WithOption(All, ModeStarred(5))}},
		{"url", "HIDDEN", []Option{WithOption(All, ModeValue("HIDDEN"))}},
		{"url", "123", []Option{WithOption(All, ModeFunc(ft))}},
		//END TEST MODES

		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/path?query#fragment", nil},

		//PASSWORD
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo@host/path?query#fragment", []Option{WithOption(Password, ModeHidden())}},
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/path?query#fragment", []Option{WithOption(Password, ModeNothing())}},
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:******@host/path?query#fragment", []Option{WithOption(Password, ModeStarred(6))}},
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:password@host/path?query#fragment", []Option{WithOption(Password, ModeValue("password"))}},
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:123@host/path?query#fragment", []Option{WithOption(Password, ModeFunc(ft))}},
		{"scheme://userinfo@host/path?query#fragment", "scheme://userinfo@host/path?query#fragment", []Option{WithOption(Password, ModeStarred(6))}},
		//END PASSWORD

		//USER
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://usr1:pwd@host/path?query#fragment", []Option{WithOption(Username, ModeValue("usr1"))}},

		//SCHEME
		{"scheme://userinfo:pwd@host/path?query#fragment", "//userinfo:pwd@host/path?query#fragment", []Option{WithOption(Scheme, ModeHidden())}},

		//PATH
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/newpath?query#fragment", []Option{WithOption(Path, ModeValue("newpath"))}},

		//QUERY. linked to a fragment. parsing with bug???
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/path?newquery#fragment", []Option{WithOption(Query, ModeValue("newquery#fragment"))}},

		//TODO: FRAGMENT. Bug?
		//{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/path?query#newfragment", []Option{WithOption(Fragment, ModeValue("newfragment"))}},

		//COMPLEX
		{
			"scheme://userinfo:pwd@host/path?query#fragment",
			"scheme://usr1:****@host/newpath?query#fragment",
			[]Option{
				WithOption(Username, ModeValue("usr1")),
				WithOption(Password, ModeStarred(4)),
				WithOption(Path, ModeValue("newpath")),
			},
		},
		//TODO: with err escape symbols
	}

	//  scheme:opaque?query#fragment
	//  scheme://userinfo@host/path?query#fragment
	for _, c := range cases {
		got := Replace(c.url, c.modes...)
		if got != c.want {
			t.Errorf("Expected %q, got %q", c.want, got)
		}
	}
}

func TestURL(t *testing.T) {
	//,
	u, err := url.Parse("scheme://userinfo:pwd@host/path?query#fragment")
	if err != nil {
		t.Fatal(err)
	}
	chk := "scheme://userinfo:****@host/path?query#fragment"
	ur := ReplaceURL(u, WithOption(Password, ModeStarred(4)))
	if got, err := url.PathUnescape(ur.String()); err != nil {
		t.Error(err)
	} else if got != chk {
		t.Errorf("Expected %q,got %q", chk, got)
	}
}

func BenchmarkReplace(b *testing.B) {
	s := "scheme://userinfo:pwd@host/path?query#fragment"
	opts := []Option{
		WithOption(Username, ModeValue("usr1")),
		WithOption(Password, ModeStarred(4)),
		WithOption(Path, ModeValue("newpath")),
		WithOption(Scheme, ModeNothing()),
		WithOption(Path, ModeNothing()),
		WithOption(Query, ModeValue("query2")),
	}
	for n := 0; n < b.N; n++ {
		Replace(s, opts...)
	}
}
